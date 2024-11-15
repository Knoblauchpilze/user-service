package service

import (
	"github.com/KnoblauchPilze/user-service/pkg/errors"
)

const (
	UserNotAuthenticated  errors.ErrorCode = 1000
	AuthenticationExpired errors.ErrorCode = 1001
)
