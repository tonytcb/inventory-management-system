package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresDatabase struct {
	instance testcontainers.Container
}

const (
	PgImage    = "postgres:14-bullseye"
	PgHost     = "host-db"
	PgUser     = "test"
	PgPassword = "test"
	PgDB       = "db_test"
	PgPort     = "5432"
)

func NewPostgresDatabase(
	t *testing.T,
	ctx context.Context,
	dockerNetwork *testcontainers.DockerNetwork,
) *PostgresDatabase {
	t.Helper()

	const timeout = 30 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if dockerNetwork == nil {
		dockerNetwork = CreateDockerNetwork(t, ctx, nil)

		t.Cleanup(func() {
			_ = dockerNetwork.Remove(ctx) //nolint:errcheck // this is for testing
		})
	}

	port := fmt.Sprintf("%s/tcp", PgPort)

	const acceptingConnectionsLog = "database system is ready to accept connections"

	req := testcontainers.ContainerRequest{
		Name:         "postgres-test-" + uuid.NewString(),
		Hostname:     PgHost,
		Image:        PgImage,
		ExposedPorts: []string{port},
		HostConfigModifier: func(config *container.HostConfig) {
			config.Mounts = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: ProjectRoot + "/internal/infra/storage/migrations",
					Target: "/docker-entrypoint-initdb.d",
				},
			}
			config.AutoRemove = true
		},
		Env: map[string]string{
			"POSTGRES_USER":     PgUser,
			"POSTGRES_PASSWORD": PgPassword,
			"POSTGRES_DB":       PgDB,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog(acceptingConnectionsLog).WithOccurrence(2).WithStartupTimeout(5*time.Second),
			wait.ForListeningPort(nat.Port(port)),
			wait.ForExec([]string{"pg_isready"}).WithStartupTimeout(5*time.Second),
		),
		Networks:       []string{dockerNetwork.Name},
		NetworkAliases: map[string][]string{dockerNetwork.Name: {PgHost}},
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	return &PostgresDatabase{
		instance: postgres,
	}
}

func (db *PostgresDatabase) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	p, err := db.instance.MappedPort(ctx, PgPort)
	require.NoError(t, err)

	return p.Int()
}

func (db *PostgresDatabase) DSN(t *testing.T) string {
	return fmt.Sprintf(
		"postgres://%s:%s@127.0.0.1:%d/%s?sslmode=disable",
		PgUser,
		PgPassword,
		db.Port(t),
		PgDB,
	)
}

func (db *PostgresDatabase) Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	require.NoError(t, db.instance.Terminate(ctx))
}
