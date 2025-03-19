package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"go-restapi-crud/utils"

	"github.com/golang-jwt/jwt/v4"
)

// JWTMiddleware is a middleware for JWT authentication
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header must be Bearer token")
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims := jwt.MapClaims{}

		// Get secret key from environment or use a default for development
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "your-secret-key" // Default secret key for development
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Extract user ID from claims
		userID, ok := claims["user_id"]
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Set user ID in request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
