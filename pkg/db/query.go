package db

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/jackc/pgx/v5"
)

func QueryOne[T any](ctx context.Context, conn Connection, sql string, arguments ...any) (T, error) {
	var out T

	connImpl, ok := conn.(*connectionImpl)
	if !ok {
		return out, errors.NewCode(UnsupportedOperation)
	}
	rows, err := connImpl.Query(ctx, sql, arguments...)
	if err != nil {
		return out, errors.WrapCode(err, QueryOneFailure)
	}

	// https://pkg.go.dev/github.com/jackc/pgx/v5#RowToStructByName
	out, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		if err == pgx.ErrNoRows {
			return out, errors.WrapCode(err, NoMatchingRows)
		} else if err == pgx.ErrTooManyRows {
			return out, errors.WrapCode(err, TooManyMatchingRows)
		}
		return out, err
	}

	return out, nil
}

func QueryAll[T any](ctx context.Context, conn Connection, sql string, arguments ...any) ([]T, error) {
	var out []T

	connImpl, ok := conn.(*connectionImpl)
	if !ok {
		return out, errors.NewCode(UnsupportedOperation)
	}
	rows, err := connImpl.Query(ctx, sql, arguments...)
	if err != nil {
		return out, errors.WrapCode(err, QueryAllFailure)
	}

	// https://pkg.go.dev/github.com/jackc/pgx/v5#RowToStructByName
	out, err = pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return out, errors.WrapCode(err, QueryAllFailure)
	}

	return out, nil
}
