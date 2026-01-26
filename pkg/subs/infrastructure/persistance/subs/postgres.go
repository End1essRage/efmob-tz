package subs

import (
	"context"
	"errors"

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

// Create сохраняет подписку
func (r *GormSubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) (uuid.UUID, error) {
	model := FromDomain(sub)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return uuid.Nil, err
	}
	return model.ID, nil // возвращаем UUID, который реально в базе
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

// Update обновляет подписку
func (r *GormSubscriptionRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	res := r.db.WithContext(ctx).Model(&SubscriptionModel{}).
		Where("id = ?", sub.ID()).
		Updates(FromDomain(sub))
	if res.RowsAffected == 0 {
		return domain.ErrSubscriptionNotFound
	}
	return res.Error
}

// Delete удаляет подписку
func (r *GormSubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&SubscriptionModel{}, "id = ?", id)
	if res.RowsAffected == 0 {
		return domain.ErrSubscriptionNotFound
	}
	return res.Error
}

// Find возвращает список подписок
func (r *GormSubscriptionRepo) Find(ctx context.Context, q domain.SubscriptionQuery, pagination p.Pagination, sorting p.Sorting) ([]*domain.Subscription, error) {
	db := r.db.WithContext(ctx).Model(&SubscriptionModel{})

	if q.UserID() != uuid.Nil {
		db = db.Where("user_id = ?", q.UserID())
	}
	if q.ServiceName() != nil {
		db = db.Where("service_name = ?", *q.ServiceName())
	}
	if q.Period() != nil {
		db = db.Where("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", q.Period().To(), q.Period().From())
	}

	if sorting.OrderBy != "" {
		db = db.Order(sorting.OrderBy + " " + string(sorting.Direction))
	}

	if pagination.Limit > 0 {
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

	if q.UserID() != uuid.Nil {
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
