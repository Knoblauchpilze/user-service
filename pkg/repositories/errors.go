package repositories

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	OptimisticLockException  errors.ErrorCode = 200
	MoreThanOneMatchingEntry errors.ErrorCode = 201
)
