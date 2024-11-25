package service

import (
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
)

const (
	UserNotAuthenticated  errors.ErrorCode = 1000
	AuthenticationExpired errors.ErrorCode = 1001
	InvalidCredentials    errors.ErrorCode = 1002

	InvalidEmail    errors.ErrorCode = 1050
	InvalidPassword errors.ErrorCode = 1051
)
