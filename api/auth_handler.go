package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/mkabdelrahman/hotel-reservation/auth"
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

	token, err := auth.GenerateAuthToken(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
