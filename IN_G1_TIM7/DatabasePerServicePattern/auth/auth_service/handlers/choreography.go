package handlers

import (
	"log"
	"net/http"

	"auth-service/choreography"
	"auth-service/database"
	"auth-service/models"
	"auth-service/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChoreographyAuthHandler - HTTP handler za koreografisani SAGA pattern.
//
// Izlaže iste endpoint-e kao i orkestrisani, ali na drugačijim putanjama:
//   POST /api/choreography/register
//   GET  /api/choreography/register/status/:sagaId
//
// Logika registracije je ista (validacija, duplikati, hash lozinke, insert),
// ali umesto da poziva Orchestrator, ovaj handler poziva ChoreographyHandler
// koji objavljuje DOGAĐAJ umesto da šalje KOMANDU.
type ChoreographyAuthHandler struct {
	db      *database.Database
	handler *choreography.ChoreographyHandler
}

func NewChoreographyAuthHandler(db *database.Database, h *choreography.ChoreographyHandler) *ChoreographyAuthHandler {
	return &ChoreographyAuthHandler{db: db, handler: h}
}

// RegisterChoreography - HTTP handler za koreografisanu registraciju
//
// POST /api/choreography/register
//
// Flow:
//  1. Validacija inputa i provera duplikata (identično orkestrisanom)
//  2. Insert u auth bazu → dobijamo newID
//  3. Objavljivanje AuthUserCreatedEvent (stakeholder sam reaguje na događaj)
//  4. HTTP 202 Accepted + sagaId
func (h *ChoreographyAuthHandler) RegisterChoreography(c *gin.Context) {
	// 1. Validacija inputa
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

	// 2. Insert u auth bazu
	var newID int64
	err = h.db.DB.QueryRow(
		`INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Username, hashedPassword, user.Email, user.Role,
	).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.Printf("[CHOREOGRAPHY][AUTH] User inserted in auth DB: %s (id=%d)", user.Username, newID)

	// 3. Objavljivanje AuthUserCreatedEvent
	// Za razliku od orkestrisanog pristupa koji šalje KOMANDU stakeholder servisu,
	// ovde objavljujemo DOGAĐAJ - stakeholder sam odlučuje da reaguje na njega.
	sagaID := uuid.New().String()
	if err := h.handler.PublishAuthUserCreated(
		sagaID, newID,
		user.Username, user.Email, user.Role,
		user.Firstname, user.Lastname,
	); err != nil {
		log.Printf("[CHOREOGRAPHY][AUTH] Failed to publish event sagaId=%s: %v", sagaID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Registration failed: could not publish AuthUserCreatedEvent. Auth entry rolled back.",
		})
		return
	}

	// 4. HTTP 202: auth DB upis je obavljen, stakeholder upis je u toku (async)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Registration initiated via choreography. AuthUserCreatedEvent published.",
		"sagaId":  sagaID,
		"pattern": "choreography",
	})
}

// RegisterChoreographyStatus - polling endpoint za status koreografisane sage
//
// GET /api/choreography/register/status/:sagaId
func (h *ChoreographyAuthHandler) RegisterChoreographyStatus(c *gin.Context) {
	sagaID := c.Param("sagaId")
	instance := h.handler.GetStatus(sagaID)
	if instance == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Saga not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sagaId":  instance.SagaID,
		"state":   instance.State,
		"user":    instance.Username,
		"created": instance.CreatedAt,
		"pattern": "choreography",
	})
}
