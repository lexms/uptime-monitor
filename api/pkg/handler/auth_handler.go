package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tou01/uptime-monitor/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	// In a real app, this would be replaced with a database repository
	users    map[string]models.User
	jwtKey   []byte
	tokenTTL time.Duration
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtKey string, tokenTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		users:    make(map[string]models.User),
		jwtKey:   []byte(jwtKey),
		tokenTTL: tokenTTL,
	}
}

// RegisterRoutes registers the auth routes to the router
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Bind JSON to request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if email is already registered
	for _, user := range h.users {
		if user.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already registered",
			})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create new user
	user := models.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      "user",
		Active:    true,
	}

	// Store user (in a real app, this would be stored in a database)
	h.users[user.ID] = user

	// Generate JWT token
	token, err := h.generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, models.AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind JSON to request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find user by email
	var foundUser models.User
	var found bool
	for _, user := range h.users {
		if user.Email == req.Email {
			foundUser = user
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.AuthResponse{
		User:  foundUser,
		Token: token,
	})
}

// generateToken generates a JWT token for the user
func (h *AuthHandler) generateToken(user models.User) (string, error) {
	// Create token claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(h.tokenTTL).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	return token.SignedString(h.jwtKey)
}
