package entities

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	UserID    string `json:"user_id"`
	ManagerID string `json:"manager_id"`
	jwt.StandardClaims
}
