package repositories

import "github.com/KnoblauchPilze/backend-toolkit/pkg/errors"

const (
	OptimisticLockException  errors.ErrorCode = 200
	MoreThanOneMatchingEntry errors.ErrorCode = 201
)
