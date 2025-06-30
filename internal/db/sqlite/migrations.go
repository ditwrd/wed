package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

func runMigrations(db *sqlx.DB) error {
	// Set up Goose with embedded migrations
	goose.SetBaseFS(migrations)

	// Set dialect for SQLite
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Run up migrations
	if err := goose.Up(db.DB, "migrations"); err != nil {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	return nil
}

func RegisterHooks(lifecycle fx.Lifecycle, db *sqlx.DB) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return runMigrations(db)
		},
		OnStop: func(ctx context.Context) error {
			err := db.Close()
			if err != nil {
				return err
			}
			return nil
		},
	})
}
