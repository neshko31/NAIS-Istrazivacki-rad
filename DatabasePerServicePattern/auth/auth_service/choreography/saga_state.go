package choreography

import (
	"sync"
	"time"
)

// SagaState predstavlja stanje jedne koreografisane SAGA instance
//
// Dijagram tranzicija (identičan kao orkestrisani, ali bez centralnog koordinatora):
//
//	STARTED
//	  │  (auth upisuje korisnika u svoju bazu i objavljuje AuthUserCreatedEvent)
//	AUTH_USER_CREATED
//	  │  (stakeholder čuje događaj, upisuje profil, objavljuje StakeholderProfileCreatedEvent)
//	COMPLETED
//
//	AUTH_USER_CREATED  →  stakeholder objavljuje AuthUserRollbackEvent
//	  │
//	COMPENSATING
//	  │  (auth čuje rolbek, briše korisnika iz auth baze)
//	COMPENSATED / FAILED
type SagaState string

const (
	SagaStateStarted         SagaState = "STARTED"
	SagaStateAuthUserCreated SagaState = "AUTH_USER_CREATED"
	SagaStateCompleted       SagaState = "COMPLETED"
	SagaStateCompensating    SagaState = "COMPENSATING"
	SagaStateCompensated     SagaState = "COMPENSATED"
	SagaStateFailed          SagaState = "FAILED"
)

// SagaInstance čuva sve podatke o jednoj koreografisanoj registracionoj transakciji
type SagaInstance struct {
	SagaID    string    `json:"sagaId"`
	State     SagaState `json:"state"`
	AuthID    int64     `json:"authId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	CreatedAt time.Time `json:"createdAt"`
}

// sagaStore čuva aktivne sage u memoriji (isti pristup kao kod orchestratora)
var sagaStore sync.Map // map[sagaId] -> *SagaInstance

func storeSaga(instance *SagaInstance) {
	sagaStore.Store(instance.SagaID, instance)
}

func loadSaga(sagaID string) *SagaInstance {
	raw, ok := sagaStore.Load(sagaID)
	if !ok {
		return nil
	}
	return raw.(*SagaInstance)
}

// NewSagaInstance kreira novu instancu sage u stanju STARTED
func NewSagaInstance(sagaID string, authID int64, username, email, role, firstname, lastname string) *SagaInstance {
	return &SagaInstance{
		SagaID:    sagaID,
		State:     SagaStateStarted,
		AuthID:    authID,
		Username:  username,
		Email:     email,
		Role:      role,
		Firstname: firstname,
		Lastname:  lastname,
		CreatedAt: time.Now(),
	}
}
