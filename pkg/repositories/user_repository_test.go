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
	assert.Nil(t, err)

	assert.Equal(t, user, actual)
	assertUserExists(t, conn, user.Id)
}

func TestIT_UserRepository_Create_WhenDuplicateName_ExpectFailure(t *testing.T) {
	repo, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	newUser := persistence.User{
		Id:        uuid.New(),
		Email:     user.Email,
		Password:  "my-password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   6,
	}

	_, err := repo.Create(context.Background(), newUser)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
	assertUserDoesNotExist(t, conn, newUser.Id)
}

func TestIT_UserRepository_Get(t *testing.T) {
	repo, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	actual, err := repo.Get(context.Background(), user.Id)
	assert.Nil(t, err)

	assert.Equal(t, user.Id, actual.Id)
	assert.Equal(t, user.Email, actual.Email)
	assert.Equal(t, user.Password, actual.Password)
	assert.Equal(t, 0, actual.Version)
}

func TestIT_UserRepository_Get_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	// Non-existent id
	id := uuid.MustParse("00000000-1111-2222-1111-000000000000")
	_, err := repo.Get(context.Background(), id)
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserRepository_GetByEmail(t *testing.T) {
	repo, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)

	actual, err := repo.GetByEmail(context.Background(), user.Email)
	assert.Nil(t, err)

	assert.Equal(t, user.Id, actual.Id)
	assert.Equal(t, user.Email, actual.Email)
	assert.Equal(t, user.Password, actual.Password)
	assert.Equal(t, 0, actual.Version)
}

func TestIT_UserRepository_GetByEmail_WhenNotFound_ExpectFailure(t *testing.T) {
	repo, _ := newTestUserRepository(t)

	_, err := repo.GetByEmail(context.Background(), "not-an-email")
	assert.True(t, errors.IsErrorWithCode(err, db.NoMatchingRows), "Actual err: %v", err)
}

func TestIT_UserRepository_List(t *testing.T) {
	repo, conn := newTestUserRepository(t)
	u1 := insertTestUser(t, conn)
	u2 := insertTestUser(t, conn)

	ids, err := repo.List(context.Background())

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(ids), 2)
	assert.Contains(t, ids, u1.Id)
	assert.Contains(t, ids, u2.Id)
}

func TestIT_UserRepository_Update(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Password = "my-new-password"

	actual, err := repo.Update(context.Background(), updatedUser)

	assert.Nil(t, err)

	assert.Equal(t, user.Id, actual.Id)
	assert.Equal(t, user.Email, actual.Email)
	assert.Equal(t, updatedUser.Password, actual.Password)
	assert.Equal(t, user.CreatedAt, actual.CreatedAt)
	assert.Equal(t, user.Version+1, actual.Version)
}

func TestIT_UserRepository_Update_WhenNameAlreadyExists_ExpectFailure(t *testing.T) {
	repo, conn := newTestUserRepository(t)
	user := insertTestUser(t, conn)
	toUpdate := insertTestUser(t, conn)

	updatedUser := toUpdate
	updatedUser.Email = user.Email

	_, err := repo.Update(context.Background(), updatedUser)

	assert.True(t, errors.IsErrorWithCode(err, pgx.UniqueConstraintViolation), "Actual err: %v", err)
}

func TestIT_UserRepository_Update_WhenVersionIsWrong_ExpectOptimisticLockException(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Password = "my-new-password"
	updatedUser.Version = user.Version + 2

	_, err := repo.Update(context.Background(), updatedUser)

	assert.True(t, errors.IsErrorWithCode(err, OptimisticLockException), "Actual err: %v", err)
}

func TestIT_UserRepository_Update_BumpsUpdatedAt(t *testing.T) {
	repo, conn := newTestUserRepository(t)

	user := insertTestUser(t, conn)

	updatedUser := user
	updatedUser.Password = "my-new-password"

	_, err := repo.Update(context.Background(), updatedUser)
	assert.Nil(t, err)

	updatedUserFromDb, err := repo.Get(context.Background(), user.Id)
	assert.Nil(t, err)
	assert.True(t, updatedUserFromDb.UpdatedAt.After(user.UpdatedAt))
}

func TestIT_UserRepository_Delete(t *testing.T) {
	repo, conn, tx := newTestUserRepositoryAndTransaction(t)

	user := insertTestUser(t, conn)

	err := repo.Delete(context.Background(), tx, user.Id)
	tx.Close(context.Background())

	assert.Nil(t, err)
	assertUserDoesNotExist(t, conn, user.Id)
}

func TestIT_UserRepository_Delete_WhenNotFound_ExpectSuccess(t *testing.T) {
	repo, conn, tx := newTestUserRepositoryAndTransaction(t)

	user := insertTestUser(t, conn)
	id := uuid.New()
	require.NotEqual(t, user.Id, id)

	err := repo.Delete(context.Background(), tx, id)
	tx.Close(context.Background())

	assert.Nil(t, err)
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
