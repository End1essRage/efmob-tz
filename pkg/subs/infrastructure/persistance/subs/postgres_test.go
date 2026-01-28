package subs

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	p "github.com/end1essrage/efmob-tz/pkg/common/persistance"
	domain "github.com/end1essrage/efmob-tz/pkg/subs/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupPostgresContainer(ctx context.Context) (tc.Container, string, error) {
	req := tc.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())
	return container, dsn, nil
}

func TestSubscriptionCommands_Integration(t *testing.T) {
	ctx := context.Background()
	container, dsn, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer container.Terminate(ctx)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	repo := NewGormSubscriptionRepo(db)

	if err := repo.Migrate(); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	sub, _ := domain.NewSubscription(uuid.Nil, uuid.New(), "service1", 100, time.Now(), nil)

	// CREATE
	uid, err := repo.Create(ctx, sub)
	assert.NoError(t, err)

	// GET — берём ID, который вернул репозиторий
	got, err := repo.GetByID(ctx, uid)
	assert.NoError(t, err)
	assert.Equal(t, uid, got.ID())

	// UPDATE — используем sub.ID() или uid, они теперь совпадают
	_ = sub.ChangePrice(150)
	err = repo.Update(ctx, sub)
	assert.NoError(t, err)
	updated, _ := repo.GetByID(ctx, uid)
	assert.Equal(t, 150, updated.Price())

	// DELETE
	err = repo.Delete(ctx, uid)
	assert.NoError(t, err)
	_, err = repo.GetByID(ctx, uid)
	assert.ErrorIs(t, err, domain.ErrSubscriptionNotFound)
}

func TestSubscriptionRepo_FindAndCalculateTotal(t *testing.T) {
	ctx := context.Background()
	container, dsn, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer container.Terminate(ctx)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	repo := NewGormSubscriptionRepo(db)
	if err := repo.Migrate(); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	userID := uuid.New()
	now := time.Now()

	// создаём несколько подписок
	subs := []*domain.Subscription{}
	for i := 0; i < 5; i++ {
		sub, _ := domain.NewSubscription(uuid.Nil, userID, fmt.Sprintf("service%d", i%2), 100*(i+1), now.AddDate(0, -i, 0), nil)
		_, err := repo.Create(ctx, sub)
		assert.NoError(t, err)
		subs = append(subs, sub)
	}

	// --- Find с фильтром ---
	query := domain.NewSubscriptionQuery(&userID, nil, nil, nil)
	pagination := &p.Pagination{Limit: 2, Offset: 0}
	sorting := &p.Sorting{OrderBy: "price", Direction: p.Descending}

	results, err := repo.Find(ctx, query, pagination, sorting)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // лимит 2
	assert.True(t, results[0].Price() > results[1].Price())

	// --- CalculateTotal ---
	total, err := repo.CalculateTotal(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
}

func TestSubscriptionRepo_ConcurrentUpdateOptimisticLocking(t *testing.T) {
	ctx := context.Background()
	container, dsn, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer container.Terminate(ctx)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewGormSubscriptionRepo(db)
	assert.NoError(t, repo.Migrate())

	// Создаем тестовую подписку
	userID := uuid.New()
	sub, err := domain.NewSubscription(uuid.Nil, userID, "test-service", 100, time.Now(), nil)
	assert.NoError(t, err)

	subscriptionID, err := repo.Create(ctx, sub)
	assert.NoError(t, err)

	// Количество конкурентных горутин
	const goroutines = 5
	var successCount int32
	var conflictCount int32
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Канал для синхронизации старта всех горутин
	start := make(chan struct{})

	// Запускаем горутины, которые пытаются обновить одну и ту же запись
	for i := 0; i < goroutines; i++ {
		go func(threadID int) {
			defer wg.Done()
			<-start // Ждем сигнала старта

			// Каждая горутина пытается обновить запись несколько раз
			for attempt := 0; attempt < 3; attempt++ {
				// Загружаем текущее состояние
				sub, err := repo.GetByID(ctx, subscriptionID)
				if err != nil {
					t.Logf("Thread %d attempt %d: failed to get subscription: %v", threadID, attempt, err)
					continue
				}

				// Меняем цену
				newPrice := sub.Price() + 1
				if err := sub.ChangePrice(newPrice); err != nil {
					t.Logf("Thread %d attempt %d: failed to change price: %v", threadID, attempt, err)
					continue
				}

				// Пытаемся сохранить
				err = repo.Update(ctx, sub)
				if err == nil {
					atomic.AddInt32(&successCount, 1)
					t.Logf("Thread %d attempt %d: успешно обновил цену на %d", threadID, attempt, newPrice)
					return // Успешно обновили, выходим из горутины
				} else if errors.Is(err, domain.ErrConcurrentModification) {
					atomic.AddInt32(&conflictCount, 1)
					t.Logf("Thread %d attempt %d: обнаружена конкурентная модификация, пробую снова", threadID, attempt)
					// Ждем немного перед повторной попыткой
					time.Sleep(time.Duration(threadID+1) * 20 * time.Millisecond)
				} else {
					t.Logf("Thread %d attempt %d: unexpected error: %v", threadID, attempt, err)
					return
				}
			}
		}(i)
	}

	// Запускаем все горутины одновременно
	close(start)
	wg.Wait()

	// Проверяем результаты
	finalSub, err := repo.GetByID(ctx, subscriptionID)
	assert.NoError(t, err)

	t.Logf("Итоговые метрики: успешных обновлений=%d, конфликтов=%d, итоговая цена=%d",
		successCount, conflictCount, finalSub.Price())

	// Ключевые проверки:
	// 1. Хотя бы одно обновление должно быть успешным
	assert.True(t, successCount >= 1, "Должно быть хотя бы одно успешное обновление")

	// 2. Итоговая цена должна быть 100 + количество успешных обновлений
	// (так как каждый раз price + 1)
	assert.Equal(t, 100+int(successCount), finalSub.Price(),
		"Итоговая цена должна быть начальная + количество успешных обновлений")

	// 3. Общее количество операций (успешных + конфликтов) должно быть >= goroutines
	totalOperations := successCount + conflictCount
	assert.True(t, totalOperations >= int32(goroutines),
		"Общее количество операций должно быть не меньше количества горутин")

	// 4. Если есть ретраи, то успешных может быть несколько
	t.Logf("Примечание: успешных обновлений %d из %d горутин. "+
		"Если есть ретраи в обработчике, это ожидаемо.", successCount, goroutines)
}

func TestSubscriptionRepo_ExplicitVersionConflict(t *testing.T) {
	ctx := context.Background()
	container, dsn, err := setupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer container.Terminate(ctx)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewGormSubscriptionRepo(db)
	assert.NoError(t, repo.Migrate())

	// Создаем подписку
	userID := uuid.New()
	sub, err := domain.NewSubscription(uuid.Nil, userID, "test-service", 100, time.Now(), nil)
	assert.NoError(t, err)

	subscriptionID, err := repo.Create(ctx, sub)
	assert.NoError(t, err)

	// Загружаем подписку в двух разных переменных
	sub1, err := repo.GetByID(ctx, subscriptionID)
	assert.NoError(t, err)

	sub2, err := repo.GetByID(ctx, subscriptionID)
	assert.NoError(t, err)

	// Меняем цену в первой копии и сохраняем
	assert.NoError(t, sub1.ChangePrice(150))
	err = repo.Update(ctx, sub1)
	assert.NoError(t, err)

	// Теперь sub2 имеет устаревшую версию
	// Пытаемся изменить вторую копию - должна быть ошибка конкурентной модификации
	assert.NoError(t, sub2.ChangePrice(200))
	err = repo.Update(ctx, sub2)
	assert.ErrorIs(t, err, domain.ErrConcurrentModification, "Должна быть ошибка конкурентной модификации")

	// Проверяем, что в БД осталась цена из первой операции
	finalSub, err := repo.GetByID(ctx, subscriptionID)
	assert.NoError(t, err)
	assert.Equal(t, 150, finalSub.Price(), "В БД должна быть цена из первой успешной операции")
}
