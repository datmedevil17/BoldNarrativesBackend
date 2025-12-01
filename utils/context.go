package utils

import (
	"errors"
	"net/http"

	"github.com/datmedevil17/backend/constants"
)

// UserContext represents user information from JWT
type UserContext struct {
	UserID string
	Email  string
}

// GetUserFromContext extracts user information from request context
func GetUserFromContext(r *http.Request) (*UserContext, error) {
	userID := r.Context().Value(constants.UserIDContextKey)
	if userID == nil {
		return nil, errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil, errors.New("invalid user ID format in context")
	}

	emailValue := r.Context().Value(constants.EmailContextKey)
	emailStr, _ := emailValue.(string) // Email is optional

	return &UserContext{
		UserID: userIDStr,
		Email:  emailStr,
	}, nil
}

// GetUserIDFromContext is a convenience function to get just the user ID
func GetUserIDFromContext(r *http.Request) (string, error) {
	userContext, err := GetUserFromContext(r)
	if err != nil {
		return "", err
	}
	return userContext.UserID, nil
}
