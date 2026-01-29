package testapp

import (
	"context"
	"testing"
	"time"

	"github.com/end1essrage/efmob-tz/pkg/subs/application/container"
	subs_repo "github.com/end1essrage/efmob-tz/pkg/subs/infrastructure/persistance/subs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestApp struct {
	Repo      *subs_repo.GormSubscriptionRepo
	Publisher *SpyEventPublisher
	Worker    *subs_repo.EventWorker
	Container *container.Container
	DB        *gorm.DB
}

// NewAppTestBuilder создаёт тестовое приложение
func NewTestApp(t *testing.T) *TestApp {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// создаем контейнер Postgres
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dsn := "host=" + host + " port=" + port.Port() + " user=test password=test dbname=testdb sslmode=disable"

	// in-memory база
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // отключаем логирование
	})

	require.NoError(t, err)

	repo := subs_repo.NewGormSubscriptionRepo(db)
	err = repo.Migrate()
	require.NoError(t, err)

	// spy publisher
	spy := &SpyEventPublisher{}

	// воркер событий с маленьким интервалом
	worker := subs_repo.NewEventWorker(db, spy, 10*time.Millisecond, 10)

	di := container.NewContainer(repo, repo, repo)

	return &TestApp{
		Repo:      repo,
		Publisher: spy,
		Worker:    worker,
		Container: di,
		DB:        db,
	}
}
