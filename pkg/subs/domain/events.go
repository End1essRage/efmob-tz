package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Event interface {
	Type() string
	MarshalJSON() ([]byte, error)
}

type SubCreatedEvent struct {
	Id     uuid.UUID
	UserID uuid.UUID
}

func (s SubCreatedEvent) Type() string {
	return "subscription_created"
}

func (s SubCreatedEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID     uuid.UUID `json:"id"`
		UserID uuid.UUID `json:"user_id"`
	}{
		ID:     s.Id,
		UserID: s.UserID,
	})
}

type SubDeletedEvent struct {
	Id uuid.UUID
}

func (s SubDeletedEvent) Type() string {
	return "subscription_deleted"
}

func (s SubDeletedEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID uuid.UUID `json:"id"`
	}{
		ID: s.Id,
	})
}
