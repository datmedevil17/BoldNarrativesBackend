package constants

// Context key types to avoid collisions
type ContextKey string

const (
	UserIDContextKey ContextKey = "userID"
	EmailContextKey  ContextKey = "email"
)
