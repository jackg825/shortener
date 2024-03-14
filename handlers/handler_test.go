package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"shortener/cache"
	"shortener/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Mock storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) SaveShortURL(shortURL, longURL, userId string) (string, error) {
	args := m.Called(shortURL, longURL, userId)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) GetLongURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

// Mock cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockCache) Set(key string, value string, ttl int) {
	m.Called(key, value, ttl)
}

func TestCreateShortURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()

	// Start PostgreSQL containerv
	pgContainer, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("../testdata", "init-url-db.sh")),
		postgres.WithDatabase("urlshortener"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	pgConnStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Start Redis container
	redisContainer, err := redis.RunContainer(ctx,
		testcontainers.WithImage("docker.io/redis:6"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	require.NoError(t, err)

	redisConnStr, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	router := gin.Default()
	store := storage.NewStorage(pgConnStr)
	cache := cache.NewCache(redisConnStr)
	router.POST("/create", CreateShortURL(store, cache))
	router.GET("/:shortUrl", RedirectShortURL(store, cache))

	requestBody := bytes.NewBufferString(`{"longUrl":"http://example.com","userId":"00000000-AAAA-BBBB-CCCC-000000000000"}`)
	req, _ := http.NewRequest("POST", "/create", requestBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	shortURLResponse := struct {
		ShortUrl string `json:"shortUrl"`
	}{}
	err = json.Unmarshal(w.Body.Bytes(), &shortURLResponse)
	assert.Equal(t, http.StatusOK, w.Code)

	require.NoError(t, err)

	req, _ = http.NewRequest("GET", "/"+shortURLResponse.ShortUrl, nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMovedPermanently, w.Code)

	defer pgContainer.Terminate(ctx)
	defer redisContainer.Terminate(ctx)
}
