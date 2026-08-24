package service

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
)

const (
	errUserNotAuthenticated  errors.ErrorCode = 1000
	errAuthenticationExpired errors.ErrorCode = 1001
	errInvalidCredentials    errors.ErrorCode = 1002

	errInvalidEmail    errors.ErrorCode = 1050
	errInvalidPassword errors.ErrorCode = 1051
)

var (
	ErrUserNotAuthenticated  = errors.FromCode(errUserNotAuthenticated)
	ErrAuthenticationExpired = errors.FromCode(errAuthenticationExpired)
	ErrInvalidCredentials    = errors.FromCode(errInvalidCredentials)

	ErrInvalidEmail    = errors.FromCode(errInvalidEmail)
	ErrInvalidPassword = errors.FromCode(errInvalidPassword)
)
