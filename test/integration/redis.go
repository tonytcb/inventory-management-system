package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Redis struct {
	instance testcontainers.Container
}

const (
	RedisImage = "redis:7.4-alpine"
	RedisPort  = "6379"
)

func NewRedis(
	t *testing.T,
	ctx context.Context,
	dockerNetwork *testcontainers.DockerNetwork,
) *Redis {
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

	port := fmt.Sprintf("%s/tcp", RedisPort)

	req := testcontainers.ContainerRequest{
		Name:         "redis-test-" + uuid.NewString(),
		Image:        RedisImage,
		ExposedPorts: []string{port},
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
		WaitingFor: wait.ForAll(
			wait.ForExec([]string{"redis-cli", "ping"}).WithStartupTimeout(5 * time.Second),
		),
		Networks: []string{dockerNetwork.Name},
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	return &Redis{
		instance: postgres,
	}
}

func (r *Redis) Port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	p, err := r.instance.MappedPort(ctx, RedisPort)
	require.NoError(t, err)

	return p.Int()
}

func (r *Redis) DSN(t *testing.T) string {
	return fmt.Sprintf("127.0.0.1:%d", r.Port(t))
}

func (r *Redis) Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	require.NoError(t, r.instance.Terminate(ctx))
}
