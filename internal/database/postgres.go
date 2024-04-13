package database

import (
	"context"
	"errors"
	"fmt"
	"redGlow/internal/config"
	"redGlow/internal/httpError"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresDB struct{
    Pool *pgxpool.Pool
}

func NewPostgresDB(cfg *config.Config, ctx context.Context, logger *zap.Logger) *PostgresDB{
	dbConnectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.PostgresUser, cfg.Postgres.PostgresPassword, cfg.Postgres.PostgresHost,
		cfg.Postgres.PostgresPort, cfg.Postgres.PostgresDatabaseName,
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

func (pdb *PostgresDB) HandleErrors(err error) httpError.HTTPError {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return httpError.NewBadRequestError(pgErr.Message)
		}
		return httpError.NewInternalServerError(pgErr.Message)
	}
	return nil
}