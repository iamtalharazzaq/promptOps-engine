package db

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/promptops/backend/pkg/db/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

var DB *bun.DB

// Init initialises the connection to Supabase/Postgres using Bun.
func Init(ctx context.Context, dsn string) (*bun.DB, error) {
	slog.Info("Connecting to database via Bun...")

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	slog.Info("Running health check on database...")
	if err := db.Ping(); err != nil {
		return nil, err
	}

	slog.Info("Running database migrations...")
	migrator := migrate.NewMigrator(db, migrations.Migrations)
	if err := migrator.Init(ctx); err != nil {
		return nil, err
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	if group.IsZero() {
		slog.Info("Database is already up to date")
	} else {
		slog.Info("Database migrated", "group", group)
	}

	DB = db
	slog.Info("Database connection established")
	return db, nil
}
