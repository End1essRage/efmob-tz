package testapp

import (
	"context"
	"fmt"
	"testing"
	"time"

	common_test "github.com/end1essrage/efmob-tz/pkg/common/testing"
	di "github.com/end1essrage/efmob-tz/pkg/subs/application/container"
	subs_repo "github.com/end1essrage/efmob-tz/pkg/subs/infrastructure/persistance/subs"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestApp struct {
	Repo        *subs_repo.GormSubscriptionRepo
	Publisher   *SpyEventPublisher
	Worker      *subs_repo.EventWorker
	Di          *di.Container
	DB          *gorm.DB
	PgContainer tc.Container
}

func (a *TestApp) Clean() {
	ctx := context.Background()
	if err := a.PgContainer.Terminate(ctx); err != nil {
		fmt.Print("")
	}
}

// NewAppTestBuilder создаёт тестовое приложение
func NewTestApp(t *testing.T) *TestApp {
	ctx := context.Background()
	container, dsn, err := common_test.SetupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

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

	di := di.NewContainer(repo, repo, repo)

	return &TestApp{
		Repo:      repo,
		Publisher: spy,
		Worker:    worker,
		Di:        di,
		DB:        db,

		PgContainer: container,
	}
}
