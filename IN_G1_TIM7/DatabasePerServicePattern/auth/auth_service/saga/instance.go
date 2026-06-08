package saga

import "time"

// SagaState predstavlja stanje jedne instance SAGA pattern-a koji obavlja registraciju
//
// Dijagram tranzicija:
//
//	STARTED
//	  │  (auth upisuje korisnika u svoju bazu)
//	AUTH_USER_CREATED
//	  │  (stakeholder upisuje korisnika u svoju bazu)
//	COMPLETED
//
//	AUTH_USER_CREATED  →  stakeholder pad
//	  │
//	COMPENSATING
//	  │  (auth briše korisnika iz svoje baze - rollback)
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

// SagaInstance čuva sve podatke o jednoj registracionoj transakciji
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
