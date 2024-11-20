package pgx

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
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

func TestIT_New_ConnectsToDatabase_WrongCredentials(t *testing.T) {
	type testCase struct {
		connStr       string
		expectedError errors.ErrorCode
	}

	testCases := []testCase{
		{
			connStr:       "postgres://test_user:comes-from-the-environment@localhost:5432/test_db",
			expectedError: AuthenticationFailed,
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			pool, err := New(context.Background(), testCase.connStr)
			require.Nil(t, err)

			err = pool.Ping(context.Background())
			require.NotNil(t, err)

			actual := AnalyzeAndWrapPgError(err)
			assert.True(t, errors.IsErrorWithCode(actual, testCase.expectedError), "Actual err: %v", err)
		})
	}

}
