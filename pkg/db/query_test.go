package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyConnection struct {
	err error
}

const sampleSqlQuery = "SELECT name FROM my_table"

func TestQueryOne_UnsupportedConnection(t *testing.T) {
	_, err := QueryOne[int](context.Background(), &dummyConnection{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation))
}

func TestQueryOne_WhenConnectionFails_ExpectError(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	sqlQuery := "SELECT name FROM my_tables"
	_, err = QueryOne[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, QueryOneFailure))

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestQueryOne_WhenNoData_ExpectFailure(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	_, err = QueryOne[element](context.Background(), conn, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, NoMatchingRows))
}

func TestQueryOne_WhenTooManyRows_ExpectFailure(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table"
	_, err = QueryOne[element](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, TooManyMatchingRows))
}

func TestQueryAll_UnsupportedConnection(t *testing.T) {
	_, err := QueryAll[int](context.Background(), &dummyConnection{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation))
}

func TestQueryAll_WhenConnectionFails_ExpectError(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	sqlQuery := "SELECT name FROM my_tables"
	_, err = QueryAll[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, QueryAllFailure))

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestQueryAll_NoData(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	out, err := QueryAll[element](context.Background(), conn, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(out)
}

func TestQueryAll_WithData(t *testing.T) {
	conn, err := New(context.Background(), dbTestConfig)
	require.Nil(t, err)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	out, err := QueryAll[element](context.Background(), conn, sqlQuery, "test-name")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Len(out, 1)
	assert.Equal("test-name", out[0].Name)
}

func (dc *dummyConnection) Close(ctx context.Context) {}

func (dc *dummyConnection) Ping(ctx context.Context) error { return dc.err }
