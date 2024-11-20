package pgx

import "github.com/KnoblauchPilze/user-service/pkg/errors"

const (
	GenericSqlError           errors.ErrorCode = 150
	ForeignKeyValidation      errors.ErrorCode = 151
	UniqueConstraintViolation errors.ErrorCode = 152
	AuthenticationFailed      errors.ErrorCode = 153
)
