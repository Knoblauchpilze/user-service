package pgx

import "github.com/KnoblauchPilze/user-service/pkg/errors"

const (
	GenericSqlError           errors.ErrorCode = 200
	ForeignKeyValidation      errors.ErrorCode = 201
	UniqueConstraintViolation errors.ErrorCode = 202
)
