package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIT_UserRepository_Create(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := persistence.User{
		Id:        uuid.New(),
		Email:     "my-email-" + uuid.New().String(),
		Password:  "my-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   6,
	}

	actual, err := repo.Create(context.Background(), user)
	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(user, actual)

	assertUserExists(t, conn, user.Id)
}

func TestIT_UserRepository_Create_WhenDuplicateName_ExpectFailure(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	user := persistence.User{
		Id:        uuid.New(),
		Email:     "user1",
		Password:  "my-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   6,
	}

	_, err := repo.Create(context.Background(), user)
	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_UserRepository_Get(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	id := uuid.MustParse("4f26321f-d0ea-46a3-83dd-6aa1c6053aaf")
	actual, err := repo.Get(context.Background(), id)
	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(id, actual.Id)
	assert.Equal("another-test-user@another-provider.com", actual.Email)
	assert.Equal("super-strong-password", actual.Password)
	assert.Equal(0, actual.Version)
}

func TestIT_UserRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), id)
	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserRepository_GetByEmail(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	actual, err := repo.GetByEmail(context.Background(), "better-test-user@mail-client.org")
	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(uuid.MustParse("00b265e6-6638-4b1b-aeac-5898c7307eb8"), actual.Id)
	assert.Equal("better-test-user@mail-client.org", actual.Email)
	assert.Equal("weakpassword", actual.Password)
	assert.Equal(0, actual.Version)
}

func TestIT_UserRepository_GetByEmail_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	_, err := repo.GetByEmail(context.Background(), "not-an-email")
	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserRepository_List(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	ids, err := repo.List(context.Background())

	assert := assert.New(t)
	assert.Nil(err)
	assert.GreaterOrEqual(len(ids), 4)
	assert.Contains(ids, uuid.MustParse("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab"))
	assert.Contains(ids, uuid.MustParse("4f26321f-d0ea-46a3-83dd-6aa1c6053aaf"))
	assert.Contains(ids, uuid.MustParse("00b265e6-6638-4b1b-aeac-5898c7307eb8"))
	assert.Contains(ids, uuid.MustParse("beb2a2dc-2a9f-48d6-b2ca-fd3b5ca3249f"))
}

func TestIT_UserRepository_Update(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Password = "my-new-password"

	actual, err := repo.Update(context.Background(), updatedUser)

	assert := assert.New(t)
	assert.Nil(err)

	assert.Equal(user.Id, actual.Id)
	assert.Equal(user.Email, actual.Email)
	assert.Equal(updatedUser.Password, actual.Password)
	assert.Equal(user.CreatedAt, actual.CreatedAt)
	assert.Equal(user.Version+1, actual.Version)
}

func TestIT_UserRepository_Update_WhenNameAlreadyExists_ExpectFailure(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Email = "i-dont-care-about-@security.de"

	_, err := repo.Update(context.Background(), updatedUser)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_UserRepository_Update_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Password = "my-new-password"
	updatedUser.Version = user.Version + 2

	_, err := repo.Update(context.Background(), updatedUser)

	assert := assert.New(t)
	assert.True(errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_UserRepository_Update_BumpsUpdatedAt(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Password = "my-new-password"

	_, err := repo.Update(context.Background(), updatedUser)
	assert := assert.New(t)
	assert.Nil(err)

	updatedUserFromDb, err := repo.Get(context.Background(), user.Id)

	assert.Nil(err)
	assert.True(updatedUserFromDb.UpdatedAt.After(user.UpdatedAt))
}

func TestIT_UserRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestUserRepositoryAndTransaction(t)

	user := insertTestUser(t, conn)

	err := repo.Delete(context.Background(), tx, user.Id)
	tx.Close(context.Background())

	assert := assert.New(t)
	assert.Nil(err)

	assertUserDoesNotExist(t, conn, user.Id)
}

func TestIT_UserRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestUserRepositoryAndTransaction(t)

	user := insertTestUser(t, conn)
	id := uuid.New()
	require.NotEqual(t, user.Id, id)

	err := repo.Delete(context.Background(), tx, id)
	tx.Close(context.Background())

	assert := assert.New(t)
	assert.Nil(err)

	assertUserExists(t, conn, user.Id)
}

func newTestUserRepository(t *testing.T) (UserRepository, db.Connection) {
	conn := newTestConnection(t)
	return NewUserRepository(conn), conn
}

func newTestUserRepositoryAndTransaction(t *testing.T) (UserRepository, db.Connection, db.Transaction) {
	conn := newTestConnection(t)
	tx, err := conn.BeginTx(context.Background())
	require.Nil(t, err)
	return NewUserRepository(conn), conn, tx
}
