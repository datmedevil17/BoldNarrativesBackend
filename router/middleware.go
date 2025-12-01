package router

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/datmedevil17/backend/constants"
	"github.com/datmedevil17/backend/utils"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// You can add logging logic here if needed
		next.ServeHTTP(w, r)
	})
}

// HealthCheck provides a health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "bold-narratives-backend",
	}
	json.NewEncoder(w).Encode(response)
}

// JWTAuthMiddleware validates JWT tokens and extracts user information
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenString, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized: " + err.Error(),
			})
			return
		}

		// Validate token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized: Invalid token",
			})
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), constants.UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, constants.EmailContextKey, claims.Email)
		r = r.WithContext(ctx)

		// Continue with the request
		next.ServeHTTP(w, r)
	})
}

// OptionalJWTMiddleware validates JWT tokens if present but doesn't require them
func OptionalJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString, err := utils.ExtractTokenFromHeader(authHeader)
			if err == nil {
				// Validate token
				claims, err := utils.ValidateJWT(tokenString)
				if err == nil {
					// Add user information to request context if token is valid
					ctx := context.WithValue(r.Context(), constants.UserIDContextKey, claims.UserID)
					ctx = context.WithValue(ctx, constants.EmailContextKey, claims.Email)
					r = r.WithContext(ctx)
				}
			}
		}

		// Continue with the request regardless of token validation
		next.ServeHTTP(w, r)
	})
}
