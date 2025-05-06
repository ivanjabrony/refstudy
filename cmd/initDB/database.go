package initDB

import (
	"context"
	"errors"
	"ivanjabrony/refstudy/cmd/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func InitDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(cfg.GetDB())
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Minute * 5

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *pgxpool.Pool, dbName string, sourceMigration string) error {
	conn, err := db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	sqlDB := stdlib.OpenDBFromPool(db)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceMigration,
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func DownMigrations(db *pgxpool.Pool, dbName string, sourceMigration string) error {
	sqlDB := stdlib.OpenDBFromPool(db)
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		sourceMigration,
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
