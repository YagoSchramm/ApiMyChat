package foundation

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgresDB opens and validates a PostgreSQL connection.
func NewPostgresDB(ctx context.Context) (*sql.DB, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = buildPostgresDSNFromEnv()
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(25)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

func buildPostgresDSNFromEnv() string {
	host := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5432")
	user := getenvDefault("DB_USER", "postgres")
	password := os.Getenv("DB_PASSWORD")
	name := getenvDefault("DB_NAME", "postgres")
	sslmode := getenvDefault("DB_SSLMODE", "disable")

	parts := []string{
		"host=" + host,
		"port=" + port,
		"user=" + user,
		"dbname=" + name,
		"sslmode=" + sslmode,
	}

	if password != "" {
		parts = append(parts, "password="+password)
	}

	return strings.Join(parts, " ")
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
