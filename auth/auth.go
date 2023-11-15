package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAuthToken(userID string) (string, error) {
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

func IsTokenNotExpired(token *jwt.Token) bool {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		expirationTime := int64(claims["expires"].(float64))
		return time.Now().Unix() < expirationTime
	}
	return false
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
