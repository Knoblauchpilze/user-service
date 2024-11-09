package pgx

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_New_InvalidConnectionString(t *testing.T) {
	pool, err := New(context.Background(), "invalid-connection-string")

	assert := assert.New(t)
	assert.Nil(pool)
	assert.NotNil(err)
	_, ok := err.(*pgconn.ParseConfigError)
	assert.True(ok)
}

func TestUnit_New_ValidConnectionString(t *testing.T) {
	const connStr = "postgres://user:password@localhost/my-db"
	pool, err := New(context.Background(), connStr)

	assert := assert.New(t)
	assert.NotNil(pool)
	assert.Nil(err)
}

func TestIT_New_ConnectsToDatabase(t *testing.T) {
	const connStr = "postgres://test_user:test_password@localhost:5432/test_db"
	pool, err := New(context.Background(), connStr)
	require.Nil(t, err)

	err = pool.Ping(context.Background())
	assert := assert.New(t)
	assert.Nil(err)
}
