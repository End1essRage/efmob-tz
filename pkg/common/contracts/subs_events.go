package contracts

const (
	UserEventsStream = "auth.user.events"

	UserRegisteredType = "user.registered"
	UserRegisteredV1   = 1
)

type UserRegisteredV1Payload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
