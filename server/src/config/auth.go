package config

import (
	"errors"
	"log"
	"time"

	"github.com/christianotieno/tasks-traker-app/server/src/entities"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken Generate a JWT token
func GenerateToken(userID int, role string, secretKey string) (string, error) {
	claims := entities.JWTClaims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}

// VerifyToken Verify and parse a JWT token
func VerifyToken(tokenString string, secretKey []byte) (*entities.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entities.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	claims, ok := token.Claims.(*entities.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
