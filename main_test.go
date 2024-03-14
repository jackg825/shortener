package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()

	log.Println("Do stuff BEFORE the tests!")
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	// Start PostgreSQL containerv
	pgContainer, _ := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join(".testdata", "init-url-db.sh")),
		postgres.WithDatabase("urlshortener"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	pgConnStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")

	// Start Redis container
	redisContainer, _ := redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	redisConnStr, _ := redisContainer.ConnectionString(ctx)

	os.Setenv("DATABASE_URL", pgConnStr)
	os.Setenv("CACHE_URL", redisConnStr)

	log.Println("Do stuff AFTER the tests!")
	defer pgContainer.Terminate(ctx)
	defer redisContainer.Terminate(ctx)

	os.Exit(exitVal)
}
func TestShortURL(t *testing.T) {
	router := gin.Default()

	requestBody := bytes.NewBufferString(`{"longUrl":"http://example.com","userId":"00000000-AAAA-BBBB-CCCC-000000000000"}`)
	req, _ := http.NewRequest("POST", "/create", requestBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	shortURLResponse := struct {
		ShortUrl string `json:"shortUrl"`
	}{}
	err := json.Unmarshal(w.Body.Bytes(), &shortURLResponse)
	require.NoError(t, err)

	// Use the shortURLResponse.ShortUrl for further testing or assertions
	fmt.Println("Generated Short URL:", shortURLResponse.ShortUrl)

	require.NoError(t, err)
}
