package subs

import (
	"context"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/common/logger"
	"github.com/end1essrage/efmob-tz/pkg/subs/application"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Type      string    `gorm:"type:text;not null"`
	Payload   []byte
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (r *GormSubscriptionRepo) CreateEvent(ctx context.Context, event domain.Event) error {
	payload, err := event.MarshalJSON()
	if err != nil {
		return err
	}

	model := EventModel{
		ID:        uuid.New(),
		Type:      event.Type(),
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	return nil
}

// EventWorker читает события из базы и публикует их
type EventWorker struct {
	db        *gorm.DB
	publisher application.EventPublisher
	interval  time.Duration
	batchSize int
}

// NewEventWorker создаёт нового воркера
func NewEventWorker(db *gorm.DB, publisher application.EventPublisher, interval time.Duration, batchSize int) *EventWorker {
	return &EventWorker{
		db:        db,
		publisher: publisher,
		interval:  interval,
		batchSize: batchSize,
	}
}

// Run запускает бесконечный цикл воркера
func (w *EventWorker) Run(ctx context.Context) {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "EventWorker",
		Func: "Run",
		Ctx:  ctx,
	})

	log.Infof("starting event worker, interval=%s, batchSize=%d", w.interval, w.batchSize)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("stopping event worker")
			return
		case <-ticker.C:
			if err := w.processBatch(ctx); err != nil {
				log.Errorf("failed to process event batch: %v", err)
			}
		}
	}
}

// processBatch читает и публикует события
func (w *EventWorker) processBatch(ctx context.Context) error {
	log := logger.Logger().WithFields(logger.LogOptions{
		Pkg:  "EventWorker",
		Func: "processBatch",
		Ctx:  ctx,
	})

	var events []EventModel

	tx := w.db.WithContext(ctx).Limit(w.batchSize).Order("created_at ASC").Find(&events)
	if tx.Error != nil {
		return tx.Error
	}

	if len(events) == 0 {
		return nil
	}

	for _, ev := range events {
		if err := w.publisher.Publish(ctx, ev.Type, ev.Payload); err != nil {
			log.Errorf("failed to publish event %s: %v", ev.ID, err)
			continue
		}

		if err := w.db.WithContext(ctx).Delete(&EventModel{}, "id = ?", ev.ID).Error; err != nil {
			log.Errorf("failed to delete published event %s: %v", ev.ID, err)
		}
	}

	return nil
}
