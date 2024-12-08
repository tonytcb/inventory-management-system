package integration

import (
	"context"
	"testing"

	"github.com/docker/docker/daemon/logger"
	"github.com/testcontainers/testcontainers-go"

	"github.com/tonytcb/inventory-management-system/internal/app/config"
)

type Dependencies struct {
	Logger   logger.Logger
	Cfg      *config.Config
	Postgres *PostgresDatabase
	Redis    *Redis
}

func NewDependencies(ctx context.Context, t *testing.T) *Dependencies {
	network := &testcontainers.DockerNetwork{
		Name: "default_test",
	}

	dependencies := &Dependencies{
		Postgres: NewPostgresDatabase(t, ctx, network),
		Redis:    NewRedis(t, ctx, network),
	}

	t.Setenv("TEST_MODE", "1")

	t.Cleanup(func() {
		//dependencies.Postgres.Close(t)
		//dependencies.Redis.Close(t)
	})

	dependencies.Cfg = &config.Config{
		Environment: "development",
		LogLevel:    "debug",
		RestAPIPort: ":8080",
		DatabaseURL: dependencies.Postgres.DSN(t),
	}

	return dependencies
}
