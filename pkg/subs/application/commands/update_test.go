package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	"github.com/end1essrage/efmob-tz/pkg/subs/application"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository реализация для тестов
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockRepository) Find(ctx context.Context, q domain.SubscriptionQuery, pagination *p.Pagination, sorting *p.Sorting) ([]*domain.Subscription, error) {
	return nil, nil
}

func (m *MockRepository) Create(ctx context.Context, sub *domain.Subscription) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

// / TestUpdateSubscriptionHandler_Validate тесты для валидации
func TestUpdateSubscriptionHandler_Validate(t *testing.T) {
	t.Parallel()

	now := time.Now()
	future := time.Date(now.Year(), now.Month()+2, 0, 0, 0, 0, 0, time.Local)
	past := time.Date(now.Year(), now.Month()-2, 0, 0, 0, 0, 0, time.Local)
	price := 100
	newPrice := 200

	// Создаем ID для подписки и пользователя
	subID := uuid.New()
	userID := uuid.New()

	// Создаем тестовую подписку с нормализованными датами
	// Важно: домен нормализует даты к началу месяца
	normalizedNow := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	normalizedFuture := time.Date(future.Year(), future.Month(), 1, 0, 0, 0, 0, time.UTC)
	normalizedPast := time.Date(past.Year(), past.Month(), 1, 0, 0, 0, 0, time.UTC)

	sub, err := domain.NewSubscription(
		subID,
		userID,
		"service",
		price,
		normalizedNow,
		&normalizedFuture,
	)
	if err != nil {
		t.Fatalf("failed to create subscription: %v", err)
	}

	// Создаем команды с нормализованными датами
	futureCmdDate := normalizedFuture
	pastCmdDate := normalizedPast

	tests := []struct {
		name     string
		sub      *domain.Subscription
		cmd      UpdateSubscriptionCommand
		expected error
	}{
		{
			name: "valid price change without date change",
			sub:  sub,
			cmd: UpdateSubscriptionCommand{
				ID:    subID, // ID обязательно
				Price: &newPrice,
			},
			expected: nil,
		},
		{
			name: "valid price change with future start date",
			sub:  sub,
			cmd: UpdateSubscriptionCommand{
				ID:        subID,
				Price:     &newPrice,
				StartDate: &futureCmdDate, // Будущая дата (после текущей)
			},
			expected: nil,
		},
		{
			name: "invalid price change with past start date",
			sub:  sub,
			cmd: UpdateSubscriptionCommand{
				ID:        subID,
				Price:     &newPrice,
				StartDate: &pastCmdDate, // Прошлая дата (раньше текущей)
			},
			expected: application.NewErrorValidationCommand("если меняется цена, то нельзя сдвигать срок начала на дату раньше"),
		},
		{
			name: "valid start date change to past",
			sub:  sub,
			cmd: UpdateSubscriptionCommand{
				StartDate: &pastCmdDate, // Прошлая дата (раньше текущей)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepository)
			handler := NewUpdateSubscriptionHandler(repo)

			err := handler.Validate(tt.sub, tt.cmd)

			if tt.expected == nil {
				assert.NoError(t, err, "Test case: %s", tt.name)
			} else {
				assert.Error(t, err, "Test case: %s", tt.name)
				assert.EqualError(t, err, tt.expected.Error(), "Test case: %s", tt.name)

				// Проверяем тип ошибки
				var valErr *application.ErrorValidationCommand
				assert.True(t, errors.As(err, &valErr),
					"Expected ErrorValidationCommand, got: %T", err)
			}
		})
	}
}
