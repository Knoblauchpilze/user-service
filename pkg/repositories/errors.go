package repositories

import "github.com/KnoblauchPilze/user-service/pkg/errors"

const (
	OptimisticLockException  errors.ErrorCode = 200
	MoreThanOneMatchingEntry errors.ErrorCode = 201
)
