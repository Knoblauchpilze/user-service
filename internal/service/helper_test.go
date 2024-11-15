package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("db_user_service", "user_service_manager", "manager_password")

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}

func insertTestUser(t *testing.T, conn db.Connection) persistence.User {
	userRepo := repositories.NewUserRepository(conn)

	id := uuid.New()
	user := persistence.User{
		Id:       id,
		Email:    fmt.Sprintf("my-user-%s", id),
		Password: "my-password",
	}
	out, err := userRepo.Create(context.Background(), user)
	require.Nil(t, err)

	return out
}

func insertApiKeyForUser(t *testing.T, userId uuid.UUID, repo repositories.ApiKeyRepository) persistence.ApiKey {
	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    userId,
		ValidUntil: time.Now().Add(3 * time.Hour),
	}

	out, err := repo.Create(context.Background(), apiKey)
	require.Nil(t, err)

	return out
}
