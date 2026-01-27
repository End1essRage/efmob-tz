package subs

import (
	"context"
	"errors"
	"math/rand"
	"time"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormSubscriptionRepo struct {
	db *gorm.DB
}

func NewGormSubscriptionRepo(db *gorm.DB) *GormSubscriptionRepo {
	return &GormSubscriptionRepo{db: db}
}

// AutoMigrate создаёт таблицу
func (r *GormSubscriptionRepo) Migrate() error {
	return r.db.AutoMigrate(&SubscriptionModel{})
}

func (r *GormSubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) (uuid.UUID, error) {
	model := FromDomain(sub)
	// Version уже должен быть 1 из домена

	var id uuid.UUID
	err := r.withRetry(ctx, func() error {
		if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
			return err
		}
		id = model.ID
		return nil
	})
	return id, err
}

// Update с оптимистичной блокировкой
func (r *GormSubscriptionRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	err := r.withRetry(ctx, func() error {
		model := FromDomain(sub)

		res := r.db.WithContext(ctx).Model(&SubscriptionModel{}).
			Where("id = ? AND version = ?", sub.ID(), sub.Version()). // Проверяем ТЕКУЩУЮ версию
			Updates(map[string]interface{}{
				"user_id":      model.UserID,
				"service_name": model.ServiceName,
				"price":        model.Price,
				"start_date":   model.StartDate,
				"end_date":     model.EndDate,
				"updated_at":   time.Now(),               // Всегда текущее время
				"version":      gorm.Expr("version + 1"), // Атомарный инкремент
			})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			// Проверяем причину
			var exists bool
			r.db.WithContext(ctx).Model(&SubscriptionModel{}).
				Select("1").
				Where("id = ?", sub.ID()).
				Limit(1).
				Find(&exists)

			if !exists {
				return domain.ErrSubscriptionNotFound
			}

			// Если запись существует, значит version не совпал
			return domain.ErrConcurrentModification
		}

		return nil
	})
	return err
}

// GetByID возвращает подписку по ID
func (r *GormSubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var m SubscriptionModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

// Delete удаляет подписку
func (r *GormSubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.withRetry(ctx, func() error {
		res := r.db.WithContext(ctx).Delete(&SubscriptionModel{}, "id = ?", id)
		if res.RowsAffected == 0 {
			return domain.ErrSubscriptionNotFound
		}
		return res.Error
	})

	return err
}

// Find возвращает список подписок
func (r *GormSubscriptionRepo) Find(ctx context.Context, q domain.SubscriptionQuery, pagination *p.Pagination, sorting *p.Sorting) ([]*domain.Subscription, error) {
	db := r.db.WithContext(ctx).Model(&SubscriptionModel{})

	if q.UserID() != nil {
		db = db.Where("user_id = ?", q.UserID())
	}
	if q.ServiceName() != nil {
		db = db.Where("service_name = ?", *q.ServiceName())
	}
	if q.Period() != nil {
		db = db.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", q.Period().To(), q.Period().From())
	}

	if sorting != nil {
		db = db.Order(sorting.OrderBy + " " + string(sorting.Direction))
	}

	if pagination != nil {
		db = db.Limit(pagination.Limit).Offset(pagination.Offset)
	}

	var models []SubscriptionModel
	if err := db.Find(&models).Error; err != nil {
		return nil, err
	}

	var result []*domain.Subscription
	for _, m := range models {
		result = append(result, m.ToDomain())
	}
	return result, nil
}

// CalculateTotal считает количество подписок
func (r *GormSubscriptionRepo) CalculateTotal(ctx context.Context, q domain.SubscriptionQuery) (int, error) {
	db := r.db.WithContext(ctx).Model(&SubscriptionModel{})

	if q.UserID() != nil {
		db = db.Where("user_id = ?", q.UserID())
	}
	if q.ServiceName() != nil {
		db = db.Where("service_name = ?", *q.ServiceName())
	}
	if q.Period() != nil {
		db = db.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", q.Period().To(), q.Period().From())
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// TODO метрика
func (r *GormSubscriptionRepo) withRetry(ctx context.Context, op func() error) error {
	const retries = 3
	var interval = 2 * time.Second

	var lastErr error
	for i := 0; i < retries; i++ {
		//jitter
		sleep := interval + time.Duration(rand.Intn(500))*time.Millisecond
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sleep):
			}
			//exponent
			interval *= 2
		}

		if err := op(); err != nil {
			if !isRetryableError(err) {
				return err
			}
			lastErr = err
			continue
		}
		return nil
	}
	return lastErr
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	// GORM может вернуть gorm.ErrInvalidTransaction, gorm.ErrInvalidSQL и др.
	// Основная цель — сетевые ошибки / connection refused
	var netErr interface{ Temporary() bool }
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	// Также можно ловить ошибки драйвера PostgreSQL
	if errors.Is(err, gorm.ErrInvalidTransaction) {
		return true
	}

	// Простейшая проверка по строке (если нет better option)
	if err.Error() == "server closed the connection unexpectedly" {
		return true
	}

	return false
}
