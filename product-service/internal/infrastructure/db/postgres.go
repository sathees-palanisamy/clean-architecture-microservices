package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/user/go-microservices/pkg/config"
	"github.com/user/go-microservices/pkg/logger"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
)

// NewConnection initializes a new database connection with pooling configuration
//
// Configuration is read from the following environment variables:
// - DB_DSN: Full connection string (overrides other DB_* vars if set)
// - DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE: Used to build DSN if DB_DSN is empty
// - DB_MAX_OPEN_CONNS: Maximum number of open connections (default: 25)
// - DB_MAX_IDLE_CONNS: Maximum number of idle connections (default: 25)
// - DB_CONN_MAX_LIFETIME_MIN: Maximum lifetime of a connection in minutes (default: 15)
// - DB_CONN_MAX_IDLE_TIME_MIN: Maximum idle time of a connection in minutes (default: 5)
func NewConnection() (*sql.DB, error) {
	// Read configuration
	dsn := config.GetEnv("DB_DSN", "")
	if dsn == "" {
		// Fallback to separate env vars if DSN is not set
		host := config.GetEnv("DB_HOST", "localhost")
		port := config.GetEnv("DB_PORT", "5432")
		user := config.GetEnv("DB_USER", "postgres")
		pass := config.GetEnv("DB_PASSWORD", "postgres")
		dbname := config.GetEnv("DB_NAME", "product_db")
		sslmode := config.GetEnv("DB_SSLMODE", "disable")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, pass, dbname, sslmode)
	}

	db, err := otelsql.Open("postgres", dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName(config.GetEnv("DB_NAME", "product_db")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pooling settings
	maxOpenConns := config.GetEnvInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := config.GetEnvInt("DB_MAX_IDLE_CONNS", 25)
	connMaxLifetimeMin := config.GetEnvInt("DB_CONN_MAX_LIFETIME_MIN", 15)
	connMaxIdleTimeMin := config.GetEnvInt("DB_CONN_MAX_IDLE_TIME_MIN", 5)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetimeMin) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(connMaxIdleTimeMin) * time.Minute)

	// Context for ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Start background stats logging
	go logDBStats(db)

	return db, nil
}

func logDBStats(db *sql.DB) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := db.Stats()
		// Using background context
		log := logger.FromContext(context.Background())
		log.Info("DB Connection Pool Stats",
			zap.Int("open_connections", stats.OpenConnections),
			zap.Int("in_use", stats.InUse),
			zap.Int("idle", stats.Idle),
			zap.Int64("wait_count", stats.WaitCount),
			zap.Duration("wait_duration", stats.WaitDuration),
			zap.Int("max_open_connections", stats.MaxOpenConnections),
		)
	}
}
