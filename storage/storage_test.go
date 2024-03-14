// This is a simplified example. You'll need to adjust it based on your actual storage implementation.
package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSaveAndGetLongURL(t *testing.T) {
	// Start PostgreSQL container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "urlshortener",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Connect to the PostgreSQL container and run your tests here.
	// Example:
	// storage := NewStorage("connectionStringToTheContainer")
	// err := storage.SaveShortURL("shortUrl", "http://example.com")
	// require.NoError(t, err)
	// longURL, err := storage.GetLongURL("shortUrl")
	// require.NoError(t, err)
	// require.Equal(t, "http://example.com", longURL)
}
