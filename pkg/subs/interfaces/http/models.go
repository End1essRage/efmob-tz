package http

import "github.com/google/uuid"

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
	UserID uuid.UUID `json:"user_id"`

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
	// required: false
	// nullable: true
	Price *int `json:"price,omitempty"`

	// Subscription start date in MM-YYYY format
	// required: false
	// nullable: true
	StartDate *string `json:"start_date,omitempty"`

	// Subscription end date in MM-YYYY format, optional
	// required: false
	// nullable: true
	EndDate *EndDateUpdate `json:"end_date,omitempty"` // Указатель для отслеживания присутствия
}

// SubscriptionQueryRequest
// swagger:model SubscriptionQueryRequest
type SubscriptionQueryRequest struct {
	// Filter by User ID (UUID), optional
	UserID *uuid.UUID `schema:"user_id,omitempty"`

	// Filter by service name, optional
	ServiceName *string `schema:"service_name,omitempty"`

	// Filter by start_date period from (MM-YYYY), optional
	StartFrom *string `schema:"start_from,omitempty"`

	// Filter by start_date period to (MM-YYYY), optional
	StartTo *string `schema:"start_to,omitempty"`

	// Filter by end_date period from (MM-YYYY), optional
	EndFrom *string `schema:"end_from,omitempty"`

	// Filter by end_date period to (MM-YYYY), optional
	EndTo *string `schema:"end_to,omitempty"`

	// Page number for pagination, optional
	Page *int `schema:"page,omitempty"`

	// Page size for pagination, optional
	PageSize *int `schema:"page_size,omitempty"`

	// Sort field, optional (e.g., "start_date" or "price")
	OrderBy *string `schema:"order_by,omitempty"`

	// Sort direction, optional: "asc" or "desc"
	Direction *string `schema:"direction,omitempty"`
}

// TotalCostRequest
// swagger:model TotalCostRequest
type TotalCostRequest struct {
	// User ID (UUID)
	UserID *uuid.UUID `schema:"user_id"`

	// Service name filter, optional
	ServiceName *string `schema:"service_name"`

	// Filter by start_date period from (MM-YYYY), optional
	StartFrom *string `schema:"start_from,omitempty"`

	// Filter by start_date period to (MM-YYYY), optional
	StartTo *string `schema:"start_to,omitempty"`

	// Filter by end_date period from (MM-YYYY), optional
	EndFrom *string `schema:"end_from,omitempty"`

	// Filter by end_date period to (MM-YYYY), optional
	EndTo *string `schema:"end_to,omitempty"`
}

// Subscription
// swagger:model Subscription
type Subscription struct {
	// ID (UUID)
	// example: bb601f22-2bf3-4721-ae6f-7636e79a0cba
	ID uuid.UUID `json:"id"`

	// User ID (UUID)
	// example: 60601fee-2bf1-4721-ae6f-7636e79a0cba
	UserID uuid.UUID `json:"user_id"`

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
