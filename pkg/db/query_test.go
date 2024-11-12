package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type dummyConnection struct {
	Connection
}

const sampleSqlQuery = "SELECT name FROM my_table"

func TestUnit_QueryOne_UnsupportedConnection(t *testing.T) {
	_, err := QueryOne[int](context.Background(), &dummyConnection{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation), "Actual err: %v", err)
}

func TestIT_QueryOne_WhenClosed_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)
	conn.Close(context.Background())

	_, err := QueryOne[int](context.Background(), conn, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, NotConnected), "Actual err: %v", err)
}

func TestIT_QueryOne_WhenConnectionFails_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)

	sqlQuery := "SELECT name FROM my_tables"
	_, err := QueryOne[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, pgx.GenericSqlError), "Actual err: %v", err)

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestIT_QueryOne_WhenNoData_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	_, err := QueryOne[element](context.Background(), conn, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, NoMatchingRows), "Actual err: %v", err)
}

func TestIT_QueryOne_WhenTooManyRows_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table"
	_, err := QueryOne[element](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, TooManyMatchingRows), "Actual err: %v", err)
}

func TestIT_QueryOne_ToStruct(t *testing.T) {
	conn := newTestConnection(t)

	type element struct {
		Id   string
		Name string
	}

	sqlQuery := "SELECT id, name FROM my_table WHERE name = 'test-name'"
	actual, err := QueryOne[element](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := element{
		Id:   "0463ed3d-bfc9-4c10-b6ee-c223bbca0fab",
		Name: "test-name",
	}
	assert.Equal(expected, actual)
}

func TestIT_QueryOne_ToString(t *testing.T) {
	conn := newTestConnection(t)

	sqlQuery := "SELECT id FROM my_table WHERE name = 'test-name'"
	actual, err := QueryOne[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab", actual)
}

func TestIT_QueryOne_ToUuid(t *testing.T) {
	conn := newTestConnection(t)

	sqlQuery := "SELECT id FROM my_table WHERE name = 'test-name'"
	actual, err := QueryOne[uuid.UUID](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := uuid.MustParse("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab")
	assert.Equal(expected, actual)
}

func TestIT_QueryAll_UnsupportedConnection(t *testing.T) {
	_, err := QueryAll[int](context.Background(), &dummyConnection{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation), "Actual err: %v", err)
}

func TestIT_QueryAll_WhenClosed_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)
	conn.Close(context.Background())

	_, err := QueryAll[int](context.Background(), conn, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, NotConnected), "Actual err: %v", err)
}

func TestIT_QueryAll_WhenConnectionFails_ExpectFailure(t *testing.T) {
	conn := newTestConnection(t)

	sqlQuery := "SELECT name FROM my_tables"
	_, err := QueryAll[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, pgx.GenericSqlError), "Actual err: %v", err)

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestIT_QueryAll_NoData(t *testing.T) {
	conn := newTestConnection(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	out, err := QueryAll[element](context.Background(), conn, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(out)
}

func TestIT_QueryAll_ToStruct(t *testing.T) {
	conn := newTestConnection(t)

	type element struct {
		Id   string
		Name string
	}

	sqlQuery := `
		SELECT
			id,
			name
		FROM
			my_table
		WHERE
			id IN (
			'0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
			'09dd5fc3-0732-4017-81e0-ffee3211d2b9'
		)`
	actual, err := QueryAll[element](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := []element{
		{
			Id:   "0463ed3d-bfc9-4c10-b6ee-c223bbca0fab",
			Name: "test-name",
		},
		{
			Id:   "09dd5fc3-0732-4017-81e0-ffee3211d2b9",
			Name: "other-name",
		},
	}
	assert.Equal(expected, actual)
}

func TestIT_QueryAll_ToString(t *testing.T) {
	conn := newTestConnection(t)

	sqlQuery := `
		SELECT
			id
		FROM
			my_table
		WHERE
			id IN (
			'0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
			'09dd5fc3-0732-4017-81e0-ffee3211d2b9'
		)`
	actual, err := QueryAll[string](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := []string{
		"0463ed3d-bfc9-4c10-b6ee-c223bbca0fab",
		"09dd5fc3-0732-4017-81e0-ffee3211d2b9",
	}
	assert.Equal(expected, actual)
}

func TestIT_QueryAll_ToUuid(t *testing.T) {
	conn := newTestConnection(t)

	sqlQuery := `
		SELECT
			id
		FROM
			my_table
		WHERE
			id IN (
			'0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
			'09dd5fc3-0732-4017-81e0-ffee3211d2b9'
		)`
	actual, err := QueryAll[uuid.UUID](context.Background(), conn, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := []uuid.UUID{
		uuid.MustParse("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab"),
		uuid.MustParse("09dd5fc3-0732-4017-81e0-ffee3211d2b9"),
	}
	assert.Equal(expected, actual)
}
