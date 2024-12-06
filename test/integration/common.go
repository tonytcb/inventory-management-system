package integration

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

var (
	_, _file, _, _ = runtime.Caller(0)                           //nolint:gochecknoglobals // to support testing
	ProjectRoot    = filepath.Join(filepath.Dir(_file), "../..") //nolint:gochecknoglobals // to support testing
)

func CreateDockerNetwork(t *testing.T, ctx context.Context, _ *string) *testcontainers.DockerNetwork {
	n, err := network.New(ctx)
	require.NoError(t, err)
	return n
}
