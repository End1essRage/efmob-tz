package subs

import (
	"context"
	"sort"
	"sync"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
)

const maxQueryCount = 100

type InMemorySubscriptionRepo struct {
	mu   sync.RWMutex
	subs map[uuid.UUID]*domain.Subscription
}

// NewInMemorySubscriptionRepo создает новый репозиторий
func NewInMemorySubscriptionRepo() *InMemorySubscriptionRepo {
	return &InMemorySubscriptionRepo{
		subs: make(map[uuid.UUID]*domain.Subscription),
	}
}

// Create сохраняет новую подписку
func (r *InMemorySubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subs[sub.ID()] = sub
	return sub.ID(), nil
}

// GetByID возвращает подписку по ID
func (r *InMemorySubscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sub, ok := r.subs[id]
	if !ok {
		return nil, domain.ErrSubscriptionNotFound
	}
	return sub, nil
}

// Update обновляет существующую подписку (ID не меняется)
func (r *InMemorySubscriptionRepo) Update(ctx context.Context, sub *domain.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := sub.ID()
	if _, ok := r.subs[id]; !ok {
		return domain.ErrSubscriptionNotFound
	}

	r.subs[id] = sub
	return nil
}

// Delete удаляет подписку по ID
func (r *InMemorySubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.subs[id]; !ok {
		return domain.ErrSubscriptionNotFound
	}
	delete(r.subs, id)
	return nil
}

// Find ищет подписки по фильтру, с пагинацией и сортировкой
func (r *InMemorySubscriptionRepo) Find(
	ctx context.Context,
	q domain.SubscriptionQuery,
	pagination *p.Pagination,
	sorting *p.Sorting,
) ([]*domain.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Subscription
	for _, sub := range r.subs {
		if q.UserID() != nil && sub.UserID() != *q.UserID() {
			continue
		}
		if q.ServiceName() != nil && sub.ServiceName() != *q.ServiceName() {
			continue
		}
		if q.Period() != nil {
			start := sub.StartDate()
			end := start
			if sub.EndDate() != nil {
				end = *sub.EndDate()
			}
			if end.Before(q.Period().From()) || start.After(q.Period().To()) {
				continue
			}
		}
		result = append(result, sub)
	}

	// сортировка
	if sorting != nil {
		sort.Slice(result, func(i, j int) bool {
			switch sorting.OrderBy {
			case "price":
				if sorting.Direction == p.Descending {
					return result[i].Price() > result[j].Price()
				}
				return result[i].Price() < result[j].Price()
			case "start_date":
				dateI := result[i].StartDate()
				dateJ := result[j].StartDate()
				if sorting.Direction == p.Descending {
					return dateI.After(dateJ)
				}
				return dateI.Before(dateJ)
			default:
				return true
			}
		})
	}

	// пагинация
	if pagination != nil {
		start := pagination.Offset
		if start > len(result) {
			start = len(result)
		}
		end := start + pagination.Limit
		if end > len(result) || pagination.Limit == 0 {
			end = len(result)
		}

		return result[start:end], nil
	}

	return result[:maxQueryCount], nil
}

func (r *InMemorySubscriptionRepo) CalculateTotal(ctx context.Context, q domain.SubscriptionQuery) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, sub := range r.subs {
		if q.UserID() != nil && sub.UserID() != *q.UserID() {
			continue
		}
		if q.ServiceName() != nil && sub.ServiceName() != *q.ServiceName() {
			continue
		}
		if q.Period() != nil {
			start := sub.StartDate()
			end := start
			if sub.EndDate() != nil {
				end = *sub.EndDate()
			}
			if end.Before(q.Period().From()) || start.After(q.Period().To()) {
				continue
			}
		}
		count++
	}

	return count, nil
}
