package db

import (
	"context"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/jackc/pgx/v5"
)

func QueryOneTx[T any](ctx context.Context, tx Transaction, sql string, arguments ...any) (T, error) {
	var out T

	txImpl, ok := tx.(*transactionImpl)
	if !ok {
		return out, errors.NewCode(UnsupportedOperation)
	}
	rows, err := txImpl.query(ctx, sql, arguments...)
	if err != nil {
		return out, errors.WrapCode(err, QueryOneFailure)
	}

	out, err = pgx.CollectExactlyOneRow(rows, getCollectorForType[T]())
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

func QueryAllTx[T any](ctx context.Context, tx Transaction, sql string, arguments ...any) ([]T, error) {
	var out []T

	txImpl, ok := tx.(*transactionImpl)
	if !ok {
		return out, errors.NewCode(UnsupportedOperation)
	}
	rows, err := txImpl.query(ctx, sql, arguments...)
	if err != nil {
		return out, errors.WrapCode(err, QueryAllFailure)
	}

	out, err = pgx.CollectRows(rows, getCollectorForType[T]())
	if err != nil {
		return out, errors.WrapCode(err, QueryAllFailure)
	}

	return out, nil
}
