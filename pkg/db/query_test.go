package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type dummyConnection struct {
	err error
}

const sampleSqlQuery = "SELECT name FROM my_table"

func TestUnit_QueryOne_UnsupportedConnection(t *testing.T) {
	_, err := QueryOne[int](context.Background(), &dummyConnection{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation))
}

func TestIT_QueryOne_WhenClosed_ExpectFailure(t *testing.T) {
	conn := NewTestConnection(t)
	conn.Close(context.Background())

	_, err := QueryOne[int](context.Background(), conn, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, QueryOneFailure))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, NotConnected))
}

func TestIT_QueryOne_WhenConnectionFails_ExpectError(t *testing.T) {
	conn := NewTestConnection(t)

	sqlQuery := "SELECT name FROM my_tables"
	_, err := QueryOne[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, QueryOneFailure))

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestIT_QueryOne_WhenNoData_ExpectFailure(t *testing.T) {
	conn := NewTestConnection(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	_, err := QueryOne[element](context.Background(), conn, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, NoMatchingRows))
}

func TestIT_QueryOne_WhenTooManyRows_ExpectFailure(t *testing.T) {
	conn := NewTestConnection(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table"
	_, err := QueryOne[element](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, TooManyMatchingRows))
}

func TestIT_QueryAll_UnsupportedConnection(t *testing.T) {
	_, err := QueryAll[int](context.Background(), &dummyConnection{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation))
}

func TestIT_QueryAll_WhenClosed_ExpectFailure(t *testing.T) {
	conn := NewTestConnection(t)
	conn.Close(context.Background())

	_, err := QueryAll[int](context.Background(), conn, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, QueryAllFailure))
	cause := errors.Unwrap(err)
	assert.True(errors.IsErrorWithCode(cause, NotConnected))
}

func TestIT_QueryAll_WhenConnectionFails_ExpectError(t *testing.T) {
	conn := NewTestConnection(t)

	sqlQuery := "SELECT name FROM my_tables"
	_, err := QueryAll[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, QueryAllFailure))

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestIT_QueryAll_NoData(t *testing.T) {
	conn := NewTestConnection(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	out, err := QueryAll[element](context.Background(), conn, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(out)
}

func TestIT_QueryAll_WithData(t *testing.T) {
	conn := NewTestConnection(t)

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
