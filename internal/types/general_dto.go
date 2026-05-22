package types



// Define custom types for keys to guarantee safety inside Gin contexts
const (
	UserIDContextKey = "ctx_user_id"
	EmailContextKey  = "ctx_user_email"
	RoleContextKey   = "ctx_user_role"
)