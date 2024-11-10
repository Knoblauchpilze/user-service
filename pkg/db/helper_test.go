package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("test_db", "test_user", "test_password")

func newTestConnection(t *testing.T) Connection {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}

func insertTestData(t *testing.T, conn Connection) (id uuid.UUID, name uuid.UUID) {
	id = uuid.New()
	name = uuid.New()
	_, err := conn.Exec(context.Background(), "INSERT INTO my_table VALUES ($1, $2)", id, name)
	require.Nil(t, err)

	return
}

func assertNameForId(t *testing.T, conn Connection, id uuid.UUID, expectedName string) {
	value, err := QueryOne[string](context.Background(), conn, "SELECT name FROM my_table WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, expectedName, value)
}

func assertIdDoesNotExist(t *testing.T, conn Connection, id uuid.UUID) {
	value, err := QueryOne[int](context.Background(), conn, "SELECT COUNT(*) FROM my_table WHERE id = $1", id)
	require.Nil(t, err)
	require.Equal(t, 0, value)
}
