package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"auth-service/database"
	"auth-service/models"
	"auth-service/saga"
	"auth-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	db              *database.Database
	jwtSecret       []byte
	tokenExpiration time.Duration
	orchestrator    *saga.RegistrationOrchestrator
}

func NewAuthHandler(db *database.Database, jwtSecret []byte, orchestrator *saga.RegistrationOrchestrator) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtSecret:       jwtSecret,
		tokenExpiration: 24 * time.Hour,
		orchestrator:    orchestrator,
	}
}

// Register - zamena za staru implementaciju; sada koristi orkestrirani SAGA pattern
//
// Flow:
//  1. Validacija inputa i provera duplikata
//  2. Insert u auth bazu -> dobijamo newID
//  3. Pokretanje orkestrirane SAGE: šalje RegisterUserCommand stakeholder servisu
//  4. HTTP 202 Accepted + sagaId (stakeholder upis je asinhron)
func (h *AuthHandler) Register(c *gin.Context) {
	//  1. Validacija inputa i provera duplikata
	var user models.UserRegister

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Provera duplikata pre upisa
	var exists bool
	err := h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", user.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking username"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already in use"})
		return
	}

	err = h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking email"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing error"})
		return
	}

	//  2. Insert u auth bazu
	var newID int64
	err = h.db.DB.QueryRow(
		`INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Username, hashedPassword, user.Email, user.Role,
	).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.Printf("[AUTH] Registered user in auth DB: %s (id=%d)", user.Username, newID)

	//  3. Pokretanje orkestrirane SAGE: šalje RegisterUserCommand stakeholder servisu
	sagaID := uuid.New().String()
	if err := h.orchestrator.StartSaga(
		sagaID, newID,
		user.Username, user.Email, user.Role,
		user.Firstname, user.Lastname,
	); err != nil {
		// Orkestrator je već pokrenuo kompenzaciju (brisanje iz auth baze)
		log.Printf("[AUTH] Saga failed to start for sagaId=%s: %v", sagaID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Registration failed: could not reach stakeholder service. Auth entry rolled back.",
		})
		return
	}

	// HTTP 202: korisnik je kreiran u auth bazi; stakeholder upis je u toku
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Registration initiated. User will be fully active once stakeholder profile is created.",
		"sagaId":  sagaID,
	})
}

// RegisterStatus - endpoint za status registracije preko SAGA pattern-a
// GET /api/register/status/:sagaId
func (h *AuthHandler) RegisterStatus(c *gin.Context) {
	sagaID := c.Param("sagaId")
	instance := h.orchestrator.GetStatus(sagaID)
	if instance == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Saga not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sagaId":  instance.SagaID,
		"state":   instance.State,
		"user":    instance.Username,
		"created": instance.CreatedAt,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var login models.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	var user models.User
	err := h.db.DB.QueryRow(
		`SELECT id, username, password, role FROM users WHERE username = $1`,
		login.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	if !utils.CheckPassword(login.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"iss":      "AuthService",
		"iat":      now.Unix(),
		"exp":      now.Add(h.tokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_in": h.tokenExpiration.Seconds(),
		"token_type": "Bearer",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
