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

func TestDeleteSubscriptionPublishesEvent(t *testing.T) {
	app := testapp.NewTestApp(t)

	userID := uuid.New()

	// создаем подписку
	sub, err := app.Container.CreateSubscriptionHandler.Handle(context.Background(), commands.CreateSubscriptionCommand{
		UserID:      userID,
		ServiceName: "Spotify",
		Price:       200,
		StartDate:   time.Now(),
	})
	require.NoError(t, err)

	// удаляем подписку
	err = app.Container.DeleteSubscriptionHandler.Handle(context.Background(), commands.DeleteSubscriptionCommand{
		ID: sub.ID(),
	})
	require.NoError(t, err)

	// запускаем воркер
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	go app.Worker.Run(ctx)

	time.Sleep(150 * time.Millisecond)

	events := app.Publisher.GetEvents()

	require.Len(t, events, 2) // Create + Delete
	require.Equal(t, domain.SubDeletedEvent{}.Type(), events[1].Topic)
}

func TestDeleteSubscriptionNotFoundDoesNotPublishEvent(t *testing.T) {
	app := testapp.NewTestApp(t)

	// пытаемся удалить несуществующую подписку
	err := app.Container.DeleteSubscriptionHandler.Handle(context.Background(), commands.DeleteSubscriptionCommand{
		ID: uuid.New(),
	})

	require.ErrorIs(t, err, domain.ErrSubscriptionNotFound)

	// запускаем воркер
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	go app.Worker.Run(ctx)

	time.Sleep(100 * time.Millisecond)

	events := app.Publisher.GetEvents()
	require.Len(t, events, 0)
}
