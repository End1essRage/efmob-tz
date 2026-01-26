package http

// TotalCostResponse
// swagger:model TotalCostResponse
type TotalCostResponse struct {
	// Total subscription cost in rubles
	// example: 1200
	Total int `json:"total"`
}

// SubscriptionCreateRequest
// swagger:model SubscriptionCreateRequest
type SubscriptionCreateRequest struct {
	// User ID (UUID)
	// required: true
	UserID string `json:"user_id"`

	// Service name
	// required: true
	ServiceName string `json:"service_name"`

	// Subscription price in rubles
	// required: true
	Price int `json:"price"`

	// Subscription start date in MM-YYYY format
	// required: true
	StartDate string `json:"start_date"`

	// Subscription end date in MM-YYYY format, optional
	// required: false
	// nullable: true
	EndDate *string `json:"end_date,omitempty"`
}

// SubscriptionUpdateRequest
// swagger:model SubscriptionUpdateRequest
type SubscriptionUpdateRequest struct {
	// Subscription price in rubles
	// required: true
	Price int `json:"price"`

	// Subscription start date in MM-YYYY format
	// required: true
	StartDate string `json:"start_date"`

	// Subscription end date in MM-YYYY format, optional
	// required: false
	// nullable: true
	EndDate *string `json:"end_date,omitempty"`
}

// SubscriptionQueryRequest
// swagger:model SubscriptionQueryRequest
type SubscriptionQueryRequest struct {
	// Filter by User ID (UUID), optional
	UserID *string `schema:"user_id"`

	// Filter by service name, optional
	ServiceName *string `schema:"service_name"`

	// Filter by start period (MM-YYYY), optional
	From *string `schema:"from"`

	// Filter by end period (MM-YYYY), optional
	To *string `schema:"to"`

	// Page number for pagination, optional
	Page int `schema:"page"`

	// Page size for pagination, optional
	PageSize int `schema:"page_size"`

	// Sort field, optional (e.g., "start_date" or "price")
	OrderBy string `schema:"order_by"`

	// Sort direction, optional: "asc" or "desc"
	Direction string `schema:"direction"`
}

// TotalCostRequest
// swagger:model TotalCostRequest
type TotalCostRequest struct {
	// User ID (UUID)
	// required: true
	UserID string `schema:"user_id"`

	// Service name filter, optional
	ServiceName *string `schema:"service_name"`

	// Period start in MM-YYYY format
	// required: true
	From string `schema:"from"`

	// Period end in MM-YYYY format
	// required: true
	To string `schema:"to"`
}

// Subscription
// swagger:model Subscription
type Subscription struct {
	// ID (UUID)
	// example: bb601f22-2bf3-4721-ae6f-7636e79a0cba
	ID string `json:"id"`

	// User ID (UUID)
	// example: 60601fee-2bf1-4721-ae6f-7636e79a0cba
	UserID string `json:"user_id"`

	// Service name
	// example: Yandex Plus
	ServiceName string `json:"service_name"`

	// Subscription price in rubles
	// example: 400
	Price int `json:"price"`

	// Subscription start date in MM-YYYY format
	// example: 07-2025
	StartDate string `json:"start_date"`

	// Subscription end date in MM-YYYY format, optional
	// required: false
	// nullable: true
	// example: 07-2026
	EndDate *string `json:"end_date,omitempty"`
}

// ErrorResponse
// swagger:response errorResponse
type ErrorResponse struct {
	// Error message
	// example: invalid request
	Error string `json:"error"`

	// Error code
	// example: VALIDATION_ERROR
	Code string `json:"code,omitempty"`
}
