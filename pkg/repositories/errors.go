package repositories

import "github.com/Knoblauchpilze/backend-toolkit/pkg/errors"

const (
	errOptimisticLockException  errors.ErrorCode = 200
	errMoreThanOneMatchingEntry errors.ErrorCode = 201
)

var (
	ErrOptimisticLockException  = errors.FromCode(errOptimisticLockException)
	ErrMoreThanOneMatchingEntry = errors.FromCode(errMoreThanOneMatchingEntry)
)
