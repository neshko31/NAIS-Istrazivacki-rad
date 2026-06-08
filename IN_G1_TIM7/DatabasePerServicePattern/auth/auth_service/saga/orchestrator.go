package saga

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RegistrationOrchestrator koordiniše dvofaznu registraciju korisnika:
//
//  1. Auth servis upisuje korisnika (username, password, email, role) u svoju Postgres bazu
//  2. Orchestrator šalje RegisterUserCommand stakeholder servisu
//  3. Stakeholder servis upisuje korisnika (profil) u svoju Postgres bazu i šalje reply
//  4. Ako stakeholder korak uspe → saga = COMPLETED
//     Ako stakeholder korak padne → Orchestrator briše korisnika iz auth baze (kompenzacija)
//
// Aktivne sage čuvaju se u memoriji (sync.Map). Za produkciju preporučujemo Redis.
type RegistrationOrchestrator struct {
	db      *sql.DB
	rmq     *RabbitMQConnection
	sagas   sync.Map // map[sagaId] -> *SagaInstance
}

func NewRegistrationOrchestrator(db *sql.DB, rmq *RabbitMQConnection) *RegistrationOrchestrator {
	o := &RegistrationOrchestrator{db: db, rmq: rmq}
	// Pokrni slušaoce reply queue-ova u pozadini
	go o.listenForReplies()
	return o
}

// ─────────────────────────────────────────────────────────────────────────────
// Korak 1: pokretanje sage (poziva se iz HTTP handler-a)
// ─────────────────────────────────────────────────────────────────────────────

// StartSaga upisuje korisnika u auth bazu, pa šalje komandu stakeholder servisu.
// Vraća sagaId koji klijent može koristiti za polling statusa.
func (o *RegistrationOrchestrator) StartSaga(
	sagaID string,
	authID int64,
	username, email, role, firstname, lastname string,
) error {
	instance := NewSagaInstance(sagaID, authID, username, email, role, firstname, lastname)
	o.sagas.Store(sagaID, instance)

	log.Printf("[ORCHESTRATION] Starting registration saga sagaId=%s user=%s", sagaID, username)

	// Korak 2: pošalji komandu stakeholder servisu
	cmd := RegisterUserCommand{
		SagaID:    sagaID,
		AuthID:    authID,
		Username:  username,
		Email:     email,
		Role:      role,
		Firstname: firstname,
		Lastname:  lastname,
	}
	body, err := json.Marshal(cmd)
	if err != nil {
		instance.State = SagaStateFailed
		return err
	}

	instance.State = SagaStateAuthUserCreated
	if err := o.rmq.Publish(RegisterUserCmdKey, body); err != nil {
		log.Printf("[ORCHESTRATION] sagaId=%s ERROR sending RegisterUserCommand: %v", sagaID, err)
		// Korisnik je upisan u auth bazu ali poruka nije stigla → kompenzuj
		o.triggerCompensation(instance)
		return err
	}

	log.Printf("[ORCHESTRATION] sagaId=%s RegisterUserCommand sent", sagaID)
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Slušanje reply-eva (pokreće se kao goroutine)
// ─────────────────────────────────────────────────────────────────────────────

func (o *RegistrationOrchestrator) listenForReplies() {
	msgs, err := o.rmq.Consume(UserRegisteredReplyQueue)
	if err != nil {
		log.Printf("[ORCHESTRATION] Failed to consume %s: %v", UserRegisteredReplyQueue, err)
		return
	}
	log.Printf("[ORCHESTRATION] Listening for replies on %s", UserRegisteredReplyQueue)

	for d := range msgs {
		o.handleUserRegisteredReply(d)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Korak 3: obrada reply-a stakeholder servisa
// ─────────────────────────────────────────────────────────────────────────────

func (o *RegistrationOrchestrator) handleUserRegisteredReply(d amqp.Delivery) {
	var reply UserRegisteredReply
	if err := json.Unmarshal(d.Body, &reply); err != nil {
		log.Printf("[ORCHESTRATION] Failed to parse UserRegisteredReply: %v", err)
		d.Nack(false, false)
		return
	}

	log.Printf("[ORCHESTRATION] Received UserRegisteredReply sagaId=%s success=%v", reply.SagaID, reply.Success)

	raw, ok := o.sagas.Load(reply.SagaID)
	if !ok {
		log.Printf("[ORCHESTRATION] Unknown sagaId=%s, ignoring", reply.SagaID)
		d.Ack(false)
		return
	}
	instance := raw.(*SagaInstance)

	if reply.Success {
		instance.State = SagaStateCompleted
		log.Printf("[ORCHESTRATION] sagaId=%s → COMPLETED: user %s registered in both databases", reply.SagaID, instance.Username)
	} else {
		log.Printf("[ORCHESTRATION] sagaId=%s stakeholder FAILED: %s, triggering compensation", reply.SagaID, reply.ErrorMessage)
		o.triggerCompensation(instance)
	}

	d.Ack(false)
}

// ─────────────────────────────────────────────────────────────────────────────
// Kompenzacija: brisanje korisnika iz auth baze
// ─────────────────────────────────────────────────────────────────────────────

func (o *RegistrationOrchestrator) triggerCompensation(instance *SagaInstance) {
	instance.State = SagaStateCompensating
	log.Printf("[ORCHESTRATION] sagaId=%s → COMPENSATING: deleting authId=%d from auth DB", instance.SagaID, instance.AuthID)

	// Direktno brišemo iz auth baze jer smo u istom servisu (Orchestrator živi u auth servisu)
	_, err := o.db.Exec("DELETE FROM users WHERE id = $1", instance.AuthID)
	if err != nil {
		log.Printf("[ORCHESTRATION] sagaId=%s CRITICAL: compensation DELETE failed: %v", instance.SagaID, err)
		instance.State = SagaStateFailed
		return
	}

	instance.State = SagaStateCompensated
	log.Printf("[ORCHESTRATION] sagaId=%s → COMPENSATED: user %s deleted from auth DB", instance.SagaID, instance.Username)
}

// GetStatus vraća trenutno stanje sage (za polling)
func (o *RegistrationOrchestrator) GetStatus(sagaID string) *SagaInstance {
	raw, ok := o.sagas.Load(sagaID)
	if !ok {
		return nil
	}
	return raw.(*SagaInstance)
}
