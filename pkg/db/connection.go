package db

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	jpgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection interface {
	Close(ctx context.Context)
	Ping(ctx context.Context) error
}

type connectionImpl struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, config Config) (Connection, error) {
	connStr := config.ToConnectionString()

	pool, err := pgx.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	conn := &connectionImpl{
		pool: pool,
	}

	return conn, err
}

func (ci *connectionImpl) Close(ctx context.Context) {
	if ci.pool != nil {
		ci.pool.Close()
		ci.pool = nil
	}
}

func (ci *connectionImpl) Ping(ctx context.Context) error {
	if ci.pool == nil {
		return errors.NewCode(NotConnected)
	}
	return ci.pool.Ping(ctx)
}

func (ci *connectionImpl) Query(ctx context.Context, sql string, arguments ...any) (jpgx.Rows, error) {
	if ci.pool == nil {
		return nil, errors.NewCode(NotConnected)
	}
	return ci.pool.Query(ctx, sql, arguments...)
}
