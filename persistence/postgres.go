package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/alexandria-oss/core/config"
	"gocloud.dev/postgres"
)

// NewPostgresPool Obtain a PostgreSQL connection pool
func NewPostgresPool(ctx context.Context, cfg *config.Kernel) (*sql.DB, func(), error) {
	db, err := postgres.Open(ctx, cfg.DBMS.URL)
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxIdleConns(10)

	cleanup := func() {
		_ = db.Close()
	}

	return db, cleanup, nil
}
