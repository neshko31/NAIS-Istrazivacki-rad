package choreography

import (
	"database/sql"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ─────────────────────────────────────────────────────────────────────────────
// ChoreographyHandler - Auth servisova uloga u koreografisanom SAGA pattern-u
//
// KLJUČNA RAZLIKA od RegistrationOrchestrator-a:
//
//   ORKESTRISANI (RegistrationOrchestrator):
//     - Zna sve korake sage unapred
//     - Šalje KOMANDE: "Stakeholder, kreiraj ovog korisnika"
//     - Prima REPLY-eve: "Evo rezultata komande"
//     - Centralna tačka kontrole - ako padne, cela saga staje
//
//   KOREOGRAFISANI (ChoreographyHandler):
//     - Ne zna ko sluša njegove događaje
//     - Objavljuje DOGAĐAJE: "Upravo sam upisao korisnika u auth bazu"
//     - Stakeholder SAM odlučuje šta da radi sa tim događajem
//     - Sluša povratne događaje i reaguje na njih
//     - Nema centralne kontrole - labavija sprežnost (loose coupling)
//
// Tok:
//   1. PublishAuthUserCreated() → objavljuje AuthUserCreatedEvent
//   2. listenForProfileCreated() → čeka StakeholderProfileCreatedEvent → COMPLETED
//   3. listenForRollback()       → čeka AuthUserRollbackEvent → kompenzacija
// ─────────────────────────────────────────────────────────────────────────────

type ChoreographyHandler struct {
	db  *sql.DB
	rmq *RabbitMQConnection
}

// NewChoreographyHandler kreira handler i pokreće listener goroutine-e u pozadini
func NewChoreographyHandler(db *sql.DB, rmq *RabbitMQConnection) *ChoreographyHandler {
	h := &ChoreographyHandler{db: db, rmq: rmq}
	// Pokreni slušaoce u pozadini (ekvivalent listenForReplies() kod orchestratora)
	go h.listenForProfileCreated()
	go h.listenForRollback()
	return h
}

// ─────────────────────────────────────────────────────────────────────────────
// Korak 1: Objavljivanje događaja (poziva se iz HTTP handler-a)
// ─────────────────────────────────────────────────────────────────────────────

// PublishAuthUserCreated čuva instancu sage i objavljuje AuthUserCreatedEvent.
// Stakeholder servis će čuti ovaj događaj i kreirati profil u svojoj bazi.
func (h *ChoreographyHandler) PublishAuthUserCreated(
	sagaID string,
	authID int64,
	username, email, role, firstname, lastname string,
) error {
	instance := NewSagaInstance(sagaID, authID, username, email, role, firstname, lastname)
	storeSaga(instance)

	log.Printf("[CHOREOGRAPHY][AUTH] Starting choreography saga sagaId=%s user=%s", sagaID, username)

	event := AuthUserCreatedEvent{
		SagaID:    sagaID,
		AuthID:    authID,
		Username:  username,
		Email:     email,
		Role:      role,
		Firstname: firstname,
		Lastname:  lastname,
	}
	body, err := json.Marshal(event)
	if err != nil {
		instance.State = SagaStateFailed
		return err
	}

	instance.State = SagaStateAuthUserCreated
	if err := h.rmq.Publish(AuthUserCreatedKey, body); err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s ERROR publishing AuthUserCreatedEvent: %v", sagaID, err)
		// Korisnik je upisan u auth bazu ali događaj nije objavljen → kompenzuj
		h.compensate(instance)
		return err
	}

	log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s AuthUserCreatedEvent published → čekamo stakeholder...", sagaID)
	return nil
}

// GetStatus vraća trenutno stanje koreografisane sage (za polling)
func (h *ChoreographyHandler) GetStatus(sagaID string) *SagaInstance {
	return loadSaga(sagaID)
}

// ─────────────────────────────────────────────────────────────────────────────
// Korak 2: Slušanje povratnih događaja od stakeholder servisa
// ─────────────────────────────────────────────────────────────────────────────

// listenForProfileCreated sluša stakeholder.profile.created.queue.
// Kada stakeholder uspešno kreira profil, objavljuje StakeholderProfileCreatedEvent.
func (h *ChoreographyHandler) listenForProfileCreated() {
	msgs, err := h.rmq.Consume(StakeholderProfileCreatedQueue)
	if err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Failed to consume %s: %v", StakeholderProfileCreatedQueue, err)
		return
	}
	log.Printf("[CHOREOGRAPHY][AUTH] Listening for StakeholderProfileCreatedEvent on %s", StakeholderProfileCreatedQueue)

	for d := range msgs {
		h.handleProfileCreated(d)
	}
}

// listenForRollback sluša auth.user.rollback.queue.
// Kada stakeholder ne uspe da kreira profil, objavljuje AuthUserRollbackEvent.
func (h *ChoreographyHandler) listenForRollback() {
	msgs, err := h.rmq.Consume(AuthUserRollbackQueue)
	if err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Failed to consume %s: %v", AuthUserRollbackQueue, err)
		return
	}
	log.Printf("[CHOREOGRAPHY][AUTH] Listening for AuthUserRollbackEvent on %s", AuthUserRollbackQueue)

	for d := range msgs {
		h.handleRollback(d)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Obrada događaja
// ─────────────────────────────────────────────────────────────────────────────

// handleProfileCreated obrađuje StakeholderProfileCreatedEvent.
// Oznaka da je stakeholder uspešno kreirao profil → saga je COMPLETED.
func (h *ChoreographyHandler) handleProfileCreated(d amqp.Delivery) {
	var event StakeholderProfileCreatedEvent
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Failed to parse StakeholderProfileCreatedEvent: %v", err)
		d.Nack(false, false)
		return
	}

	log.Printf("[CHOREOGRAPHY][AUTH] Received StakeholderProfileCreatedEvent sagaId=%s success=%v",
		event.SagaID, event.Success)

	instance := loadSaga(event.SagaID)
	if instance == nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Unknown sagaId=%s, ignoring", event.SagaID)
		d.Ack(false)
		return
	}

	instance.State = SagaStateCompleted
	log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s → COMPLETED: user %s registered in both databases",
		event.SagaID, instance.Username)
	d.Ack(false)
}

// handleRollback obrađuje AuthUserRollbackEvent.
// Stakeholder nije uspeo → brišemo korisnika iz auth baze (kompenzacija).
func (h *ChoreographyHandler) handleRollback(d amqp.Delivery) {
	var event AuthUserRollbackEvent
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Failed to parse AuthUserRollbackEvent: %v", err)
		d.Nack(false, false)
		return
	}

	log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s received rollback event: %s",
		event.SagaID, event.ErrorMessage)

	instance := loadSaga(event.SagaID)
	if instance == nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Unknown sagaId=%s, ignoring rollback", event.SagaID)
		d.Ack(false)
		return
	}

	h.compensate(instance)
	d.Ack(false)
}

// ─────────────────────────────────────────────────────────────────────────────
// Kompenzacija: brisanje korisnika iz auth baze
// (identična logika kao kod orchestratora)
// ─────────────────────────────────────────────────────────────────────────────

func (h *ChoreographyHandler) compensate(instance *SagaInstance) {
	instance.State = SagaStateCompensating
	log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s → COMPENSATING: deleting authId=%d from auth DB",
		instance.SagaID, instance.AuthID)

	_, err := h.db.Exec("DELETE FROM users WHERE id = $1", instance.AuthID)
	if err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s CRITICAL: compensation DELETE failed: %v",
			instance.SagaID, err)
		instance.State = SagaStateFailed
		return
	}

	instance.State = SagaStateCompensated
	log.Printf("[CHOREOGRAPHY][AUTH] sagaId=%s → COMPENSATED: user %s deleted from auth DB",
		instance.SagaID, instance.Username)
}
