package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a JWT token using user ID and email
func GenerateJWT(userID uint, email string) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key from .env
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}
