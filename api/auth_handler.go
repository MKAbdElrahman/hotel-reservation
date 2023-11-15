package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store db.UserStore
}

func NewAuthHandler(store db.UserStore) *AuthHandler {
	return &AuthHandler{store: store}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *gin.Context) {
	var authParams AuthParams

	// Bind the JSON request body to the AuthParams struct
	if err := c.BindJSON(&authParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.store.GetUserByEmail(c, authParams.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the provided password with the hashed password stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(authParams.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ok, err := types.IsValidPassword(user.EncryptedPassword, authParams.Password)

	if err != nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateAuthToken(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateAuthToken(userID string) (string, error) {
	// Set standard claims
	claims := jwt.MapClaims{
		"id":      userID,                                  // Set subject claim to the user ID
		"expires": time.Now().Add(time.Second * 60).Unix(), // Token expiration time (e.g., 24 hours)
		"issued":  time.Now().Unix(),                       // Token issue time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET")) // Replace with your actual secret key

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}
