package testing

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresContainer wraps a PostgreSQL test container
type PostgresContainer struct {
	Container testcontainers.Container
	DB        *gorm.DB
	DSN       string
}

// SetupPostgres creates a PostgreSQL test container
func SetupPostgres(t *testing.T) *PostgresContainer {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		container.Terminate(ctx)
	})

	return &PostgresContainer{
		Container: container,
		DB:        db,
		DSN:       dsn,
	}
}

// RedisContainer wraps a Redis test container
type RedisContainer struct {
	Container testcontainers.Container
	Client    *redis.Client
	Addr      string
}

// SetupRedis creates a Redis test container
func SetupRedis(t *testing.T) *RedisContainer {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	addr := fmt.Sprintf("%s:%s", host, port.Port())
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
		container.Terminate(ctx)
	})

	return &RedisContainer{
		Container: container,
		Client:    client,
		Addr:      addr,
	}
}

// TruncateTables truncates all tables in the database
func TruncateTables(t *testing.T, db *gorm.DB, tables ...string) {
	t.Helper()

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

// FlushRedis flushes all Redis data
func FlushRedis(t *testing.T, client *redis.Client) {
	t.Helper()

	if err := client.FlushAll(context.Background()).Err(); err != nil {
		t.Fatalf("Failed to flush Redis: %v", err)
	}
}

// RunInTransaction runs a function in a database transaction and rolls back
func RunInTransaction(t *testing.T, db *gorm.DB, fn func(tx *gorm.DB)) {
	t.Helper()

	tx := db.Begin()
	defer tx.Rollback()

	fn(tx)
}

// AssertNoError is a helper to assert no error occurred
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// AssertError is a helper to assert an error occurred
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
}

// AssertEqual is a helper to assert two values are equal
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %v, want %v", got, want)
	}
}

// AssertNotEqual is a helper to assert two values are not equal
func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Fatalf("Got %v, want not equal", got)
	}
}

// AssertTrue is a helper to assert a condition is true
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Fatalf("Assertion failed: %s", message)
	}
}

// AssertFalse is a helper to assert a condition is false
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Fatalf("Assertion failed: %s", message)
	}
}

// WaitFor waits for a condition to be true with timeout
func WaitFor(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Timeout waiting for condition")
}
