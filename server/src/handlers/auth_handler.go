package handlers

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file", err)
		}

		secret := os.Getenv("SECRET")

		// Verify and parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Provide the same JWT signing key used during token generation
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Check if the token is valid and has not expired
		if !token.Valid {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Access the user ID from the JWT claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Pass the user ID to the next handler
		ctx := context.WithValue(r.Context(), "userID", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
