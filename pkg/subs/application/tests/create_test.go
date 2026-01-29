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
	// Создаем тестовое приложение
	app := testapp.NewTestApp(t)
	userID := uuid.New()

	// Создаем подписку через handler
	sub, err := app.Container.CreateSubscriptionHandler.Handle(context.Background(), commands.CreateSubscriptionCommand{
		UserID:      userID,
		ServiceName: "Netflix",
		Price:       100,
		StartDate:   time.Now(),
		EndDate:     nil,
	})
	require.NoError(t, err)
	require.NotNil(t, sub)

	// Запускаем воркер в фоне, чтобы обработал событие
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	go app.Worker.Run(ctx)

	// Ждем немного, чтобы воркер успел обработать событие
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что событие опубликовано
	events := app.Publisher.GetEvents()
	require.Len(t, events, 1)
	require.Equal(t, domain.SubCreatedEvent{}.Type(), events[0].Topic)
}
