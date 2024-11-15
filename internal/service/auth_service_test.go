package service

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockApiKeyRepository struct {
	repositories.ApiKeyRepository

	apiKey persistence.ApiKey
	err    error
}

func TestUnit_AuthService_Authenticate_WhenKeyDoesNotExist_ExpectFailure(t *testing.T) {
	repo := &mockApiKeyRepository{
		err: errors.NewCode(db.NoMatchingRows),
	}
	auth := newTestAuthRepository(repo)

	_, err := auth.Authenticate(context.Background(), uuid.New())

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, UserNotAuthenticated), "Actual err: %v", err)
}

func TestUnit_AuthService_Authenticate_WhenKeyExpired_ExpectFailure(t *testing.T) {
	dateInThePast, _ := time.Parse(time.RFC3339, "2024-11-15T01:00:00Z")
	repo := &mockApiKeyRepository{
		apiKey: persistence.ApiKey{
			ValidUntil: dateInThePast,
		},
	}
	auth := newTestAuthRepository(repo)

	_, err := auth.Authenticate(context.Background(), uuid.New())

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, AuthenticationExpired), "Actual err: %v", err)
}

func TestIT_AuthService_Authenticate_WhenAuthenticated_ExpectSuccess(t *testing.T) {
	conn := newTestConnection(t)
	user := insertTestUser(t, conn)
	repos := repositories.Repositories{
		ApiKey: repositories.NewApiKeyRepository(conn),
	}
	apiKey := insertApiKeyForUser(t, user.Id, repos.ApiKey)

	auth := NewAuthService(repos)
	_, err := auth.Authenticate(context.Background(), apiKey.Key)

	assert := assert.New(t)
	assert.Nil(err)
}

func (m *mockApiKeyRepository) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	return m.apiKey, m.err
}

func newTestAuthRepository(apiKeyRepo repositories.ApiKeyRepository) AuthService {
	repos := repositories.Repositories{
		ApiKey: apiKeyRepo,
	}
	return NewAuthService(repos)
}
