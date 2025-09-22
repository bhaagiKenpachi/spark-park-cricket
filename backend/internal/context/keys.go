package context

// Context key types for avoiding collisions
type contextKey string

const (
	UserKey          contextKey = "user"
	UserIDKey        contextKey = "user_id"
	UserEmailKey     contextKey = "user_email"
	AuthenticatedKey contextKey = "authenticated"
	IsAdminKey       contextKey = "is_admin"
)
