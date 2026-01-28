package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSubscription_OK(t *testing.T) {
	userID := uuid.New()
	start := time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC)

	sub, err := NewSubscription(
		uuid.Nil,
		userID,
		"Yandex Plus",
		400,
		start,
		nil,
	)

	require.NoError(t, err)
	require.Equal(t, userID, sub.UserID())
	require.Equal(t, 400, sub.Price())
	require.Equal(t, time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC), sub.StartDate())
	require.Nil(t, sub.EndDate())
}

func TestNewSubscription_InvalidPrice(t *testing.T) {
	_, err := NewSubscription(
		uuid.Nil,
		uuid.New(),
		"Netflix",
		0,
		time.Now(),
		nil,
	)

	require.ErrorIs(t, err, ErrInvalidPrice)
}

func TestSubscription_ChangePrice(t *testing.T) {
	sub := mustSubscription(t)

	err := sub.ChangePrice(999)

	require.NoError(t, err)
	require.Equal(t, 999, sub.Price())
}

func TestSubscription_ChangePrice_Invalid(t *testing.T) {
	sub := mustSubscription(t)

	err := sub.ChangePrice(-10)

	require.ErrorIs(t, err, ErrInvalidPrice)
}

func TestSubscription_IsActive(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	sub, _ := NewSubscription(
		uuid.Nil,
		uuid.New(),
		"Test",
		100,
		start,
		&end,
	)

	require.False(t, sub.IsActive(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)))
	require.True(t, sub.IsActive(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)))
	require.False(t, sub.IsActive(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func mustSubscription(t *testing.T) *Subscription {
	t.Helper()

	sub, err := NewSubscription(
		uuid.Nil,
		uuid.New(),
		"Test",
		100,
		time.Now(),
		nil,
	)

	require.NoError(t, err)
	return sub
}
