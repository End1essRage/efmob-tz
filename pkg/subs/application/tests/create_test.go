//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/subs/application/commands"
	"github.com/end1essrage/efmob-tz/pkg/subs/application/tests/testapp"
	"github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateSubscriptionPublishesEvent(t *testing.T) {
	app := testapp.NewTestApp(t)
	userID := uuid.New()

	sub, err := app.Di.CreateSubscriptionHandler.Handle(context.Background(), commands.CreateSubscriptionCommand{
		UserID:      userID,
		ServiceName: "Netflix",
		Price:       100,
		StartDate:   time.Now(),
		EndDate:     nil,
	})
	require.NoError(t, err)
	require.NotNil(t, sub)

	// Запускаем воркер на 1 секунду
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go app.Worker.Run(ctx)

	// Ждем, пока событие появится, максимум 1 сек
	var events []testapp.PublishedEvent
	for start := time.Now(); time.Since(start) < time.Second; {
		events = app.Publisher.GetEvents()
		if len(events) > 0 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	require.Len(t, events, 1)
	require.Equal(t, domain.SubCreatedEvent{}.Type(), events[0].Topic)
}
