package choreography

// ─────────────────────────────────────────────────────────────────────────────
// Koreografisani SAGA pattern - Events (Auth servis)
//
// Razlika od orkestrisanog pristupa:
//   ORKESTRISANI  → Orchestrator šalje KOMANDE ("uradi ovo") i prima REPLY-eve
//   KOREOGRAFISANI → Svaki servis objavljuje DOGAĐAJE ("ovo se desilo") i
//                    ostali servisi sami odlučuju šta da urade
//
// Tok registracije (koreografija):
//   1. Auth servis  → objavljuje AuthUserCreatedEvent    (upisao sam korisnika)
//   2. Stakeholder  → čuje događaj, upisuje profil
//   3a. Stakeholder → objavljuje StakeholderProfileCreatedEvent (uspelo)
//   3b. Stakeholder → objavljuje AuthUserRollbackEvent         (palo, rolbek!)
//   4a. Auth servis → čuje uspeh → COMPLETED
//   4b. Auth servis → čuje rolbek → briše iz auth baze → COMPENSATED
// ─────────────────────────────────────────────────────────────────────────────

// AuthUserCreatedEvent - događaj koji auth servis objavljuje nakon upisa korisnika u auth bazu.
// Stakeholder servis ga čuje i kreira profil u svojoj bazi.
type AuthUserCreatedEvent struct {
	SagaID    string `json:"sagaId"`
	AuthID    int64  `json:"authId"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// StakeholderProfileCreatedEvent - događaj koji stakeholder servis objavljuje
// kada uspešno kreira profil korisnika u stakeholder bazi.
// Auth servis ga čuje i označava sagu kao COMPLETED.
type StakeholderProfileCreatedEvent struct {
	SagaID  string `json:"sagaId"`
	Success bool   `json:"success"`
}

// AuthUserRollbackEvent - događaj koji stakeholder servis objavljuje
// kada NE uspe da kreira profil korisnika.
// Auth servis ga čuje i briše korisnika iz auth baze (kompenzacija).
type AuthUserRollbackEvent struct {
	SagaID       string `json:"sagaId"`
	ErrorMessage string `json:"errorMessage"`
}
