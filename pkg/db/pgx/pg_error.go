package pgx

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	foreignKeyViolation = "23503"
	uniqueValidation    = "23505"
)

func AnalyzeAndWrapPgError(err error) error {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return err
	}

	switch pgErr.Code {
	case foreignKeyViolation:
		return errors.WrapCode(err, ForeignKeyValidation)
	case uniqueValidation:
		return errors.WrapCode(err, UniqueConstraintViolation)
	}

	return errors.WrapCode(err, GenericSqlError)
}
