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

func TestSubscriptionCommands_CRUD(t *testing.T) {
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
	query := domain.NewSubscriptionQuery(&userID, nil, nil, nil, nil)
	pagination := p.Pagination{Limit: 2, Offset: 0}
	sorting := &p.Sorting{OrderBy: "price", Direction: p.Descending}

	results, err := repo.Find(ctx, query, pagination, sorting)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // лимит 2
	assert.True(t, results[0].Price() > results[1].Price())

	// --- CalculateTotal ---
	total, err := repo.CalculateTotalCost(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, 1500, total)
}

func TestSubscriptionRepo_Queries(t *testing.T) {
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
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// -------------------------
	// Фиксированные подписки
	// -------------------------
	// Подписки с end_date = NULL
	// ID | start_date | end_date
	// 0  | 2024-01-01 | NULL
	// 1  | 2023-12-01 | NULL
	// 2  | 2023-11-01 | NULL
	// 3  | 2023-10-01 | NULL
	// 4  | 2023-09-01 | NULL
	// 5  | 2023-08-01 | NULL
	nullStarts := []time.Time{
		base, base.AddDate(0, -1, 0), base.AddDate(0, -2, 0),
		base.AddDate(0, -3, 0), base.AddDate(0, -4, 0), base.AddDate(0, -5, 0),
	}
	subs := []*domain.Subscription{}
	for _, s := range nullStarts {
		sub, err := domain.NewSubscription(uuid.Nil, userID, "service", 100, s, nil)
		assert.NoError(t, err)
		_, err = repo.Create(ctx, sub)
		assert.NoError(t, err)
		subs = append(subs, sub)
	}

	// Подписки с end_date != NULL
	endedSubs := []struct {
		start time.Time
		end   time.Time
	}{
		{base.AddDate(0, -2, 0), base.AddDate(0, -1, 0)}, // 2023-11 → 2023-12
		{base.AddDate(0, -3, 0), base.AddDate(0, -2, 0)}, // 2023-10 → 2023-11
		{base.AddDate(0, -4, 0), base.AddDate(0, -3, 0)}, // 2023-09 → 2023-09
		{base.AddDate(0, -1, 0), base.AddDate(0, 0, 0)},  // 2023-11 → 2024-01
		{base.AddDate(0, -3, 0), base.AddDate(0, -1, 0)}, // 2023-10 → 2023-12
		{base.AddDate(0, 1, 0), base.AddDate(0, 2, 0)},   // 2024-02 → 2024-03
	}
	for _, s := range endedSubs {
		sub, err := domain.NewSubscription(uuid.Nil, userID, "service", 100, s.start, &s.end)
		assert.NoError(t, err)
		_, err = repo.Create(ctx, sub)
		assert.NoError(t, err)
		subs = append(subs, sub)
	}

	// -------------------------
	// Фильтры
	// -------------------------
	startFrom := base.AddDate(0, -2, 0) // 2023-11-01
	startTo := base.AddDate(0, -1, 0)   // 2023-12-01

	endFrom := base.AddDate(0, -1, 0) // 2023-12-01
	endTo := base.AddDate(0, 1, 0)    // 2024-02-01

	// startFrom(2023-11) : nullStarts-3 endedSubs-3
	// startTo(2023-12) : nullStarts-5 endedSubs-5
	// startFrom(2023-11) - startTo(2023-12) : nullStarts-2 endedSubs-2

	// endFrom(2023-12) : nullStarts-6 endedSubs-4
	// endTo(2024-02) : nullStarts-0 endedSubs-5
	// endFrom(2023-12) - endTo(2024-02) : nullStarts-0 endedSubs-3

	// -------------------------
	// Тестовые кейсы
	// -------------------------
	tests := []struct {
		name      string
		query     domain.SubscriptionQuery
		wantCount int
	}{
		{
			name:      "start_from default:(NULL included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, nil), nil, nil),
			wantCount: 6, // start >= 2023-11: nullStarts:3 endedSubs:3
		},
		{
			name:      "start_from (NULL NOT included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, nil), nil, ptrBool(false)),
			wantCount: 3, // start >= 2023-11: nullStarts:0 endedSubs:3
		},
		{
			name:      "start_to default:(NULL included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(nil, &startTo), nil, nil),
			wantCount: 10, // start <= 2023-12: nullStarts:5 endedSubs:5
		},
		{
			name:      "start_to default:(NULL NOT included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(nil, &startTo), nil, ptrBool(false)),
			wantCount: 5, // start <= 2023-12: nullStarts:0 endedSubs:5
		},
		{
			name:      "start_from + start_to",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, &startTo), nil, nil),
			wantCount: 4, // start между 2023-11 и 2023-12: nullStarts:2 endedSubs:2
		},
		{
			name:      "end_to default:(NULL NOT included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(nil, &endTo), nil),
			wantCount: 5, // все что закончатся до 2024-02: nullStarts:0 endedSubs:5
		},
		{
			name:      "end_to (NULL included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(nil, &endTo), ptrBool(true)),
			wantCount: 11, // все что закончатся до 2024-02 или не указано: nullStarts:6 endedSubs:5
		},
		{
			name:      "end_to (NULL NOT included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(nil, &endTo), ptrBool(false)),
			wantCount: 5, // все что закончатся до 2024-02: nullStarts:0 endedSubs:5
		},
		{
			name:      "end_from default:(NULL included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, nil), nil),
			wantCount: 10, // все что закончатся после 2023-12 или не указано: nullStarts:6 endedSubs:4
		},
		{
			name:      "end_from (NULL included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, nil), ptrBool(true)),
			wantCount: 10, // все что закончатся после 2023-12 или не указано: nullStarts:6 endedSubs:4
		},
		{
			name:      "end_from (NULL NOT included)",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, nil), ptrBool(false)),
			wantCount: 4, // все что закончатся после 2023-12 или когда-то: nullStarts:0 endedSubs:4
		},
		{
			name:      "start_from + start_to + end_from",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, &startTo), mustPeriod(&endFrom, nil), nil),
			wantCount: 4, // start между 2023-11 и 2023-12 И end после 2023-12 или не указано: nullStarts:2 endedSubs:2
		},
		{
			name:      "start_from + start_to + end_to",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, &startTo), mustPeriod(nil, &endTo), nil),
			wantCount: 2, // start между 2023-11 и 2023-12 И end до 2024-02: nullStarts:0 endedSubs:2
		},
		{
			name:      "start_from + end_from + end_to",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, nil), mustPeriod(&endFrom, &endTo), nil),
			wantCount: 2, // start >= 2023-11 И end между 2023-12 и 2024-02: nullStarts:0 endedSubs:2
		},
		{
			name:      "start_to + end_from + end_to",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(nil, &startTo), mustPeriod(&endFrom, &endTo), nil),
			wantCount: 3, // start <= 2023-12 И end между 2023-12 и 2024-02: nullStarts:0 endedSubs:3
		},
		{
			name:      "start_from + start_to + end_from + end_to",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, &startTo), mustPeriod(&endFrom, &endTo), nil),
			wantCount: 2, // start между 2023-11-01 и 2023-12-01 И end между 2023-12 и 2024-02: nullStarts:0 endedSubs:2
		},
		{
			name:      "end_from + end_to + EndIsNil=null",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, &endTo), ptrBool(false)),
			wantCount: 3, // end между 2023-12 и 2024-02:  nullStarts:0 endedSubs:3
		},
		{
			name:      "end_from + end_to EndIsNil=true",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, &endTo), ptrBool(true)),
			wantCount: 9, //end между 2023-12 и 2024-02 или не заполнен:  nullStarts:6 endedSubs:3
		},
		{
			name:      "end_from + end_to + EndIsNil=false",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, &endTo), ptrBool(false)),
			wantCount: 3, // end между 2023-12 и 2024-02:  nullStarts:0 endedSubs:3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.Find(ctx, tt.query, p.DefaultPagination(), nil)
			assert.NoError(t, err)
			assert.Len(t, results, tt.wantCount, "не совпадает количество подписок")
		})
	}
}

func TestSubscriptionRepo_BoundaryMonthFilters(t *testing.T) {
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

	// Подписки, даты нормализованы, день = 1
	subs := []struct {
		start time.Time
		end   *time.Time
	}{
		{time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC), nil},                                                   // start 2023-11-01, end NULL
		{time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC), ptrTime(time.Date(2024, 01, 1, 0, 0, 0, 0, time.UTC))}, // start 2023-12-01, end 2024-01-01
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil},                                                    // start 2024-01-01, end NULL
	}

	for _, s := range subs {
		sub, err := domain.NewSubscription(uuid.Nil, userID, "service", 100, s.start, s.end)
		assert.NoError(t, err)
		_, err = repo.Create(ctx, sub)
		assert.NoError(t, err)
	}

	// Фильтры по границам месяца
	startFrom := time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC)
	startTo := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)

	endFrom := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	endTo := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		query     domain.SubscriptionQuery
		wantCount int
	}{
		{
			name:      "start_from boundary",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, nil), nil, nil),
			wantCount: 3, // все подписки с start >= 2023-11-01
		},
		{
			name:      "start_to boundary",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(nil, &startTo), nil, nil),
			wantCount: 2, // start <= 2023-12-01
		},
		{
			name:      "end_from boundary",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(&endFrom, nil), ptrBool(true)),
			wantCount: 3, // end >= 2023-12-01 или NULL
		},
		{
			name:      "end_to boundary",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(nil, &endTo), nil),
			wantCount: 0, // end <= 2023-12-01
		},
		{
			name:      "end_to boundary",
			query:     domain.NewSubscriptionQuery(&userID, nil, nil, mustPeriod(nil, ptrTime(endTo.AddDate(0, 1, 0))), nil),
			wantCount: 1, // end <= 2024-01-01
		},
		{
			name:      "start + end boundary",
			query:     domain.NewSubscriptionQuery(&userID, nil, mustPeriod(&startFrom, &startTo), mustPeriod(&endFrom, &endTo), ptrBool(true)),
			wantCount: 1, // только подписка 2023-12-01 с end 2023-12-01
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.Find(ctx, tt.query, p.DefaultPagination(), nil)
			assert.NoError(t, err)
			assert.Len(t, results, tt.wantCount, "не совпадает количество подписок")
		})
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
func ptrBool(v bool) *bool { return &v }

func mustPeriod(from, to *time.Time) *domain.Period {
	p, err := domain.NewPeriod(from, to)
	if err != nil {
		panic(err)
	}
	return p
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
				} else if errors.Is(err, ErrConcurrentModification) {
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
	assert.ErrorIs(t, err, ErrConcurrentModification, "Должна быть ошибка конкурентной модификации")

	// Проверяем, что в БД осталась цена из первой операции
	finalSub, err := repo.GetByID(ctx, subscriptionID)
	assert.NoError(t, err)
	assert.Equal(t, 150, finalSub.Price(), "В БД должна быть цена из первой успешной операции")
}
