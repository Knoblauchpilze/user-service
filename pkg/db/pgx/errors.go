package pgx

import "github.com/KnoblauchPilze/user-service/pkg/errors"

const (
	ForeignKeyValidation      errors.ErrorCode = 200
	UniqueConstraintViolation errors.ErrorCode = 201
)
