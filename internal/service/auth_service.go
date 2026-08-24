package service

import (
	"context"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/user-service/pkg/communication"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type AuthService interface {
	Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error)
}

type authServiceImpl struct {
	apiKeyRepo repositories.ApiKeyRepository
}

func NewAuthService(repos repositories.Repositories) AuthService {
	return &authServiceImpl{
		apiKeyRepo: repos.ApiKey,
	}
}

func (s *authServiceImpl) Authenticate(ctx context.Context, apiKey uuid.UUID) (communication.AuthorizationDtoResponse, error) {
	var out communication.AuthorizationDtoResponse

	key, err := s.apiKeyRepo.GetForKey(ctx, apiKey)
	if err != nil {
		if err == db.ErrNoMatchingRows {
			return out, ErrUserNotAuthenticated
		}

		return out, err
	}

	if key.ValidUntil.Before(time.Now()) {
		return out, ErrAuthenticationExpired
	}

	return out, nil
}
