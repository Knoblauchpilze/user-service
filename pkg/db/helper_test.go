package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
	"github.com/stretchr/testify/require"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("test_db", "test_user", "test_password")

func NewTestConnection(t *testing.T) Connection {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)
	return conn
}
