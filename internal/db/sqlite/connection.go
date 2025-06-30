package sqlite

import (
	"context"
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func NewDBConnection() (*sqlx.DB, error) {
	const initSQL = `
	PRAGMA journal_mode = WAL;
	PRAGMA synchronous = NORMAL;
	PRAGMA temp_store = MEMORY;
	PRAGMA busy_timeout = 5000;
	PRAGMA automatic_index = true;
	PRAGMA foreign_keys = ON;
	PRAGMA analysis_limit = 1000;
	PRAGMA trusted_schema = OFF;
`
	sqlite.RegisterConnectionHook(func(conn sqlite.ExecQuerierContext, _ string) error {
		_, err := conn.ExecContext(context.Background(), initSQL, nil)
		return err
	})

	dsn := viper.GetString("app.db.dsn")
	if dsn == "" {
		dsn = "file:/data/local.db"
	}
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Pool/health settings
	db.SetMaxOpenConns(viper.GetInt("app.db.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("app.db.max_idle_conns"))
	if v := viper.GetDuration("app.db.conn_max_lifetime"); v > 0 {
		db.SetConnMaxLifetime(v)
	}
	if v := viper.GetDuration("app.db.conn_max_idle_time"); v > 0 {
		db.SetConnMaxIdleTime(v)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
