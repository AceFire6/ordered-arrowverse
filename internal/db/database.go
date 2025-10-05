package db

import (
	"context"
	"fmt"

	pgxzerolog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type DBTypeName string

func CreateDBPool(log *zerolog.Logger, cfg *Config, typesToRegister []DBTypeName) (*pgxpool.Pool, error) {
	parsedConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	logLevel, err := tracelog.LogLevelFromString(cfg.LogLevel.String())
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	parsedConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzerolog.NewLogger(*log),
		LogLevel: logLevel,
	}
	// This runs after each connection is established to register the types with the driver
	parsedConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, typeName := range typesToRegister {
			dbType, typeErr := conn.LoadType(ctx, string(typeName))
			if typeErr != nil {
				return fmt.Errorf("load type %s: %w", typeName, typeErr)
			}
			conn.TypeMap().RegisterType(dbType)
		}

		return nil
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to DB: %w", err)
	}

	return dbPool, nil
}

func CloseDBPool(dbPool *pgxpool.Pool, logger *zerolog.Logger) {
	logger.Info().Msg("Closing DB connection")
	dbPool.Close()
}
