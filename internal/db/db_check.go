package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dbCheck struct {
	dbPool *pgxpool.Pool
}

func NewDBCheck(dbPool *pgxpool.Pool) *dbCheck {
	return &dbCheck{dbPool: dbPool}
}

func (dbc *dbCheck) RunCheck(ctx context.Context) error {
	err := dbc.dbPool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}

func (dbc *dbCheck) Name() string {
	return "database"
}
