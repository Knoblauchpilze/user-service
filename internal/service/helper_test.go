package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/Knoblauchpilze/user-service/pkg/persistence"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("db_user_service", "user_service_manager", "manager_password")

func newTestConnection(t *testing.T) *db.Connection {
	conn, err := db.New(t.Context(), dbTestConfig)
	require.Nil(t, err)
	return conn
}

func insertTestUser(t *testing.T, conn *db.Connection) persistence.User {
	t.Helper()
	repo := repositories.NewUserRepository(conn)

	id := uuid.New()
	user := persistence.User{
		Id:        id,
		Email:     fmt.Sprintf("my-user-%s", id),
		Password:  "my-password",
		CreatedAt: time.Now(),
	}
	out, err := repo.Create(t.Context(), user)
	require.Nil(t, err)

	assertUserExists(t, conn, out.Id)

	return out
}

func insertApiKeyForUser(t *testing.T, conn *db.Connection, userId uuid.UUID) persistence.ApiKey {
	t.Helper()
	return insertApiKeyForUserWithValidity(t, conn, userId, time.Now().Add(3*time.Hour))
}

func insertApiKeyForUserWithValidity(t *testing.T, conn *db.Connection, userId uuid.UUID, validity time.Time) persistence.ApiKey {
	t.Helper()
	repo := repositories.NewApiKeyRepository(conn)

	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    userId,
		ValidUntil: validity,
	}

	out, err := repo.Create(t.Context(), apiKey)
	require.Nil(t, err)

	assertApiKeyExists(t, conn, out.Id)

	return out
}

func assertApiKeyExists(t *testing.T, conn *db.Connection, id uuid.UUID) {
	t.Helper()
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, "SELECT id FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertApiKeyExistsByKey(t *testing.T, conn *db.Connection, key uuid.UUID) {
	t.Helper()
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, "SELECT key FROM api_key WHERE key = $1", key)
	require.Nil(t, err)
	require.Equal(t, key, value)
}

func assertApiKeyDoesNotExist(t *testing.T, conn *db.Connection, id uuid.UUID) {
	t.Helper()
	value, err := db.QueryOne[int](t.Context(), conn, "SELECT COUNT(id) FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func assertUserExists(t *testing.T, conn *db.Connection, id uuid.UUID) {
	t.Helper()
	value, err := db.QueryOne[uuid.UUID](t.Context(), conn, "SELECT id FROM api_user WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertUserDoesNotExist(t *testing.T, conn *db.Connection, id uuid.UUID) {
	t.Helper()
	value, err := db.QueryOne[int](t.Context(), conn, "SELECT COUNT(id) FROM api_user WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}
