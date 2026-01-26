package subs

import (
	"context"
	"fmt"
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
	query := domain.NewSubscriptionQuery(userID, nil, nil)
	pagination := p.Pagination{Limit: 2, Offset: 0}
	sorting := p.Sorting{OrderBy: "price", Direction: p.Descending}

	results, err := repo.Find(ctx, query, pagination, sorting)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // лимит 2
	assert.True(t, results[0].Price() > results[1].Price())

	// --- CalculateTotal ---
	total, err := repo.CalculateTotal(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
}
