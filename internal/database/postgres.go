package database

import (
	"context"
	"fmt"
	"redGlow/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresDB struct{
    Pool *pgxpool.Pool
}

func NewPostgresDB(cfg *config.Config, ctx context.Context, logger *zap.Logger) *PostgresDB{
	dbConnectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresDB.PostgresUser, cfg.PostgresDB.PostgresPassword, cfg.PostgresDB.PostgresHost,
		cfg.PostgresDB.PostgresPort, cfg.PostgresDB.PostgresDatabaseName,
	)

    pool, err := pgxpool.New(ctx, dbConnectString)
	if err != nil {
		logger.Fatal("Error to connect to postgresql")
	}

    return &PostgresDB{
        Pool: pool,
    }
}

func (pdb *PostgresDB) Connect(ctx context.Context) (*pgxpool.Conn,error){
	return pdb.Pool.Acquire(ctx)
}