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
	repo := repositories.NewUserRepository(conn)

	id := uuid.New()
	user := persistence.User{
		Id:        id,
		Email:     fmt.Sprintf("my-user-%s", id),
		Password:  "my-password",
		CreatedAt: time.Now(),
	}
	out, err := repo.Create(context.Background(), user)
	require.Nil(t, err)

	assertUserExists(t, conn, out.Id)

	return out
}

func insertApiKeyForUser(t *testing.T, conn db.Connection, userId uuid.UUID) persistence.ApiKey {
	return insertApiKeyForUserWithValidity(t, conn, userId, time.Now().Add(3*time.Hour))
}

func insertApiKeyForUserWithValidity(t *testing.T, conn db.Connection, userId uuid.UUID, validity time.Time) persistence.ApiKey {
	repo := repositories.NewApiKeyRepository(conn)

	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    userId,
		ValidUntil: validity,
	}

	out, err := repo.Create(context.Background(), apiKey)
	require.Nil(t, err)

	assertApiKeyExists(t, conn, out.Id)

	return out
}

func assertApiKeyExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT id FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertApiKeyExistsByKey(t *testing.T, conn db.Connection, key uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT key FROM api_key WHERE key = $1", key)
	require.Nil(t, err)
	require.Equal(t, key, value)
}

func assertApiKeyDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[int](context.Background(), conn, "SELECT COUNT(id) FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertUserExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT id FROM api_user WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertUserDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[int](context.Background(), conn, "SELECT COUNT(id) FROM api_user WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}
