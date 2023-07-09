package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/christianotieno/tasks-traker-app/server/src/config"

	"github.com/joho/godotenv"
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
		claims, err := config.VerifyToken(tokenString, []byte(secret))
		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Pass the user ID and role to the next handler
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
