package http

/*
// RegisterRequest represents registration request
// swagger:parameters registerUser
type RegisterRequest struct {
	// Required: true
	// Example: user@example.com
	Email string `json:"email"`

	// Required: true
	// Minimum: 6
	// Example: password123
	Password string `json:"password"`

	// Required: true
	// Example: John Doe
	Name string `json:"name"`
}

// RegisterResponse represents registration response
// swagger:response registerResponse
type RegisterResponse struct {
	// Response message
	// example: user created
	Message string `json:"message"`

	// User ID
	// example: 1
	UserID int `json:"user_id,omitempty"`
}
*/

// ErrorResponse represents error response
// swagger:response errorResponse
type ErrorResponse struct {

	// Error message
	// example: invalid request
	Error string `json:"error"`

	// Error code
	// example: VALIDATION_ERROR
	Code string `json:"code,omitempty"`
}
