package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"auth-service/database"
	"auth-service/models"
	"auth-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	db              *database.Database
	jwtSecret       []byte
	tokenExpiration time.Duration
}

func NewAuthHandler(db *database.Database, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		db:              db,
		jwtSecret:       jwtSecret,
		tokenExpiration: 24 * time.Hour,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
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

	var exists bool
	err := h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM registry WHERE username = $1)", user.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking username"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already in use"})
		return
	}

	err = h.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM registry WHERE email = $1)", user.Email).Scan(&exists)
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

	tx, err := h.db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed to start"})
		return
	}

	var newID int64
	err = tx.QueryRow(
		`INSERT INTO registry (username, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Username, hashedPassword, user.Email, user.Role,
	).Scan(&newID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create registry entry"})
		return
	}

	_, err = tx.Exec(
		`INSERT INTO users (registry_id, username, email, role, first_name, last_name, is_blocked)
		 VALUES ($1, $2, $3, $4, $5, $6, false)`,
		newID, user.Username, user.Email, user.Role, user.Firstname, user.Lastname,
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed to commit"})
		return
	}

	log.Printf("Registered user: %s (id=%d)", user.Username, newID)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var login models.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}

	var user models.User
	err := h.db.DB.QueryRow(
		`SELECT id, username, password, role FROM registry WHERE username = $1`,
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
