package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/db/postgresql"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("db_user_service", "user_service_manager", "manager_password")

func newTestConnection(t *testing.T) db.Connection {
	conn, err := db.New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
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

func insertTestUser(t *testing.T, conn db.Connection) persistence.User {
	someTime := time.Date(2024, 11, 12, 17, 55, 30, 0, time.UTC)

	user := persistence.User{
		Id:        uuid.New(),
		Email:     "my-email-" + uuid.New().String(),
		Password:  "my-password",
		CreatedAt: someTime,
	}
	updatedAt, err := db.QueryOne[time.Time](context.Background(), conn, "INSERT INTO api_user (id, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING updated_at", user.Id, user.Email, user.Password, user.CreatedAt)
	require.Nil(t, err)

	user.UpdatedAt = updatedAt

	return user
}
