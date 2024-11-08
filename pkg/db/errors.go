package db

import "github.com/KnoblauchPilze/user-service/pkg/errors"

const (
	NotConnected         errors.ErrorCode = 100
	UnsupportedOperation errors.ErrorCode = 101

	NoMatchingRows      errors.ErrorCode = 110
	TooManyMatchingRows errors.ErrorCode = 111
	QueryOneFailure     errors.ErrorCode = 112
	QueryAllFailure     errors.ErrorCode = 113
)
