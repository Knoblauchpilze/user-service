package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_ApiKeyRepository_Create(t *testing.T) {
	repo, conn := newTestApiKeyRepository(t)

	user := insertTestUser(t, conn)

	apiKey := persistence.ApiKey{
		Id:      uuid.New(),
		Key:     uuid.New(),
		ApiUser: user.Id,

		ValidUntil: time.Date(2024, 11, 12, 18, 32, 20, 0, time.UTC),
	}

	actual, err := repo.Create(context.Background(), apiKey)
	assert.Nil(t, err)

	assert.Equal(t, apiKey, actual)
	assertApiKeyExists(t, conn, apiKey.Id)
}

func TestIT_ApiKeyRepository_Create_WhenDuplicateForUser_ExpectKeyIsReturned(t *testing.T) {
	repo, _ := newTestApiKeyRepository(t)

	apiKey := persistence.ApiKey{
		Id:      uuid.New(),
		Key:     uuid.New(),
		ApiUser: uuid.MustParse("4f26321f-d0ea-46a3-83dd-6aa1c6053aaf"),

		ValidUntil: time.Date(2024, 11, 12, 18, 34, 40, 0, time.UTC),
	}

	actual, err := repo.Create(context.Background(), apiKey)

	assert.Nil(t, err)
	assert.Equal(t, uuid.MustParse("fd8136c4-c584-4bbf-a390-53d5c2548fb8"), actual.Id)
	assert.Equal(t, apiKey.ApiUser, actual.ApiUser)
	assert.Equal(t, uuid.MustParse("2da3e9ec-7299-473a-be0f-d722d870f51a"), actual.Key)
}

func TestIT_ApiKeyRepository_Create_WhenDuplicateForUser_ExpectValidityExtended(t *testing.T) {
	repo, _ := newTestApiKeyRepository(t)

	old, err := repo.Get(context.Background(), uuid.MustParse("fd8136c4-c584-4bbf-a390-53d5c2548fb8"))
	require.Nil(t, err)

	apiKey := persistence.ApiKey{
		Id:      old.Id,
		Key:     old.Key,
		ApiUser: old.ApiUser,

		ValidUntil: time.Date(2024, 11, 12, 18, 34, 40, 0, time.UTC),
	}

	actual, err := repo.Create(context.Background(), apiKey)
	require.Nil(t, err)

	updated, err := repo.Get(context.Background(), old.Id)
	require.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, apiKey.ValidUntil, actual.ValidUntil)
	assert.Equal(t, apiKey.ValidUntil, updated.ValidUntil.UTC())
}

func TestIT_ApiKeyRepository_Get(t *testing.T) {
	repo, conn := newTestApiKeyRepository(t)

	_, apiKey := insertTestApiKey(t, conn)

	actual, err := repo.Get(context.Background(), apiKey.Id)
	assert.Nil(t, err)

	actualUtc := actual
	actualUtc.ValidUntil = actual.ValidUntil.UTC()
	assert.Equal(t, apiKey, actualUtc)
}

func TestIT_ApiKeyRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestApiKeyRepository(t)

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_ApiKeyRepository_GetForKey(t *testing.T) {
	repo, conn := newTestApiKeyRepository(t)

	_, apiKey := insertTestApiKey(t, conn)

	actual, err := repo.GetForKey(context.Background(), apiKey.Key)
	assert.Nil(t, err)

	actualUtc := actual
	actualUtc.ValidUntil = actual.ValidUntil.UTC()
	assert.Equal(t, apiKey, actualUtc)
}

func TestIT_ApiKeyRepository_GetForKey_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestApiKeyRepository(t)

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.GetForKey(context.Background(), id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_ApiKeyRepository_GetForUser(t *testing.T) {
	repo, conn := newTestApiKeyRepository(t)

	_, apiKey := insertTestApiKey(t, conn)

	actual, err := repo.GetForUser(context.Background(), apiKey.ApiUser)
	assert.Nil(t, err)

	actualUtc := actual
	actualUtc.ValidUntil = actual.ValidUntil.UTC()
	assert.Equal(t, apiKey, actualUtc)
}

func TestIT_ApiKeyRepository_GetForUser_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestApiKeyRepository(t)

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.GetForUser(context.Background(), id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_ApiKeyRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestApiKeyRepositoryAndTransaction(t)

	user, apiKey := insertTestApiKey(t, conn)

	err := repo.DeleteForUser(context.Background(), tx, user.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertApiKeyDoesNotExist(t, conn, apiKey.Id)
}

func TestIT_ApiKeyRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestApiKeyRepositoryAndTransaction(t)

	user, apiKey := insertTestApiKey(t, conn)
	id := uuid.New()
	require.NotEqual(t, user.Id, id)

	err := repo.DeleteForUser(context.Background(), tx, id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertApiKeyExists(t, conn, apiKey.Id)
}

func newTestApiKeyRepository(t *testing.T) (ApiKeyRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewApiKeyRepository(conn), conn
}

func newTestApiKeyRepositoryAndTransaction(t *testing.T) (ApiKeyRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewApiKeyRepository(conn), conn, tx
}

func assertApiKeyExists(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[uuid.UUID](context.Background(), conn, "SELECT id FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, id, value)
}

func assertApiKeyDoesNotExist(t *testing.T, conn db.Connection, id uuid.UUID) {
	value, err := db.QueryOne[int](context.Background(), conn, "SELECT COUNT(id) FROM api_key WHERE id = $1", id)
	require.Nil(t, err)
	require.Zero(t, value)
}

func insertTestApiKey(t *testing.T, conn db.Connection) (persistence.User, persistence.ApiKey) {
	user := insertTestUser(t, conn)

	someTime := time.Date(2024, 11, 12, 18, 49, 35, 0, time.UTC)

	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    user.Id,
		ValidUntil: someTime,
	}
	_, err := conn.Exec(context.Background(), "INSERT INTO api_key (id, key, api_user, valid_until) VALUES ($1, $2, $3, $4)", apiKey.Id, apiKey.Key, apiKey.ApiUser, apiKey.ValidUntil)
	require.Nil(t, err)

	return user, apiKey
}
