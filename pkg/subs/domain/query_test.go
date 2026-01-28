package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPeriod_OK(t *testing.T) {
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	p, err := NewPeriod(&from, &to)

	require.NoError(t, err)
	require.Equal(t, from, *p.From())
	require.Equal(t, to, *p.To())
}
