package entities

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
