package database

import (
	"context"
	"email-verification-service/internal/pkg/config"
	"email-verification-service/migration"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Database) *Database {
	return &Database{Pool: initDatabase(ctx, cfg)}
}

func initDatabase(ctx context.Context, cfg config.Database) *pgxpool.Pool {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
	slog.With("connection string", connectionString).Debug("Connecting to database")
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil
	}
	err = migration.MigrateDatabase(connectionString)
	if err != nil {
		log.Fatalf("Unable to migrate database: %v", err)
		return nil
	}
	return pool
}
