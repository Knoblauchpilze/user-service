package repositories

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user persistence.User) (persistence.User, error)
	Get(ctx context.Context, id uuid.UUID) (persistence.User, error)
	GetByEmail(ctx context.Context, email string) (persistence.User, error)
	List(ctx context.Context) ([]uuid.UUID, error)
	Update(ctx context.Context, user persistence.User) (persistence.User, error)
	Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error
}

type userRepositoryImpl struct {
	conn db.Connection
}

func NewUserRepository(conn db.Connection) UserRepository {
	return &userRepositoryImpl{
		conn: conn,
	}
}

const createUserSqlTemplate = `
INSERT INTO api_user (id, email, password, created_at)
	VALUES($1, $2, $3, $4)
	RETURNING updated_at`

func (r *userRepositoryImpl) Create(ctx context.Context, user persistence.User) (persistence.User, error) {
	updatedAt, err := db.QueryOne[time.Time](ctx, r.conn, createUserSqlTemplate, user.Id, user.Email, user.Password, user.CreatedAt)
	user.UpdatedAt = updatedAt
	return user, err
}

const getUserSqlTemplate = `
SELECT
	id, email, password, created_at, updated_at, version
FROM
	api_user
WHERE
	id = $1`

func (r *userRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	return db.QueryOne[persistence.User](ctx, r.conn, getUserSqlTemplate, id)
}

const getUserByEmailSqlTemplate = `
SELECT
	id, email, password, created_at, updated_at, version
FROM
	api_user
WHERE
	email = $1`

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (persistence.User, error) {
	return db.QueryOne[persistence.User](ctx, r.conn, getUserByEmailSqlTemplate, email)
}

const listUserSqlTemplate = `
SELECT
	id
FROM
	api_user`

func (r *userRepositoryImpl) List(ctx context.Context) ([]uuid.UUID, error) {
	return db.QueryAll[uuid.UUID](ctx, r.conn, listUserSqlTemplate)
}

const updateUserSqlTemplate = `
UPDATE
	api_user
SET
	email = $1,
	password = $2,
	version = $3
WHERE
	id = $4
	AND version = $5
RETURNING
	updated_at`

func (r *userRepositoryImpl) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	version := user.Version + 1

	updatedAt, err := db.QueryOne[time.Time](ctx, r.conn, updateUserSqlTemplate, user.Email, user.Password, version, user.Id, user.Version)
	if err != nil {
		if errors.IsErrorWithCode(err, db.NoMatchingRows) {
			return user, errors.NewCode(OptimisticLockException)
		}
		return user, err
	}

	user.Version = version
	user.UpdatedAt = updatedAt

	return user, nil
}

const deleteUserSqlTemplate = `
DELETE FROM
	api_user
WHERE
	id = $1`

func (r *userRepositoryImpl) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	_, err := tx.Exec(ctx, deleteUserSqlTemplate, id)
	return err
}
