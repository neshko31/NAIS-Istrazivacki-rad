package saga

// Komande (auth servis ih šalje, stakeholder servis ih prima)

// RegisterUserCommand - komanda koju Orchestrator šalje stakeholder servisu da kreira korisnika u stakeholder bazi
type RegisterUserCommand struct {
	SagaID    string `json:"sagaId"`
	AuthID    int64  `json:"authId"`    // ID iz auth baze (registry_id u stakeholder bazi)
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// DeleteUserCommand - kompenzaciona komanda; briše korisnika iz auth baze ako je stakeholder korak pao
type DeleteUserCommand struct {
	SagaID string `json:"sagaId"`
	AuthID int64  `json:"authId"`
}

// Odogovori

// UserRegisteredReply - odgovor stakeholder servisa na RegisterUserCommand
type UserRegisteredReply struct {
	SagaID       string `json:"sagaId"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// UserDeletedReply - odgovor auth servisa na sopstvenu kompenzaciju
type UserDeletedReply struct {
	SagaID  string `json:"sagaId"`
	Success bool   `json:"success"`
}
