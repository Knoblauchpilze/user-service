package db

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/db/pgx"
	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type dummyTransaction struct {
	Transaction
}

func TestUnit_QueryOneTx_UnsupportedConnection(t *testing.T) {
	_, err := QueryOneTx[int](context.Background(), &dummyTransaction{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation))
}

func TestIT_QueryOneTx_WhenCommitted_ExpectFailure(t *testing.T) {
	_, tx := newTestTransaction(t)
	tx.Close(context.Background())

	_, err := QueryOneTx[int](context.Background(), tx, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, AlreadyCommitted))
}

func TestIT_QueryOneTx_WhenConnectionFails_ExpectFailure(t *testing.T) {
	_, tx := newTestTransaction(t)

	sqlQuery := "SELECT name FROM my_tables"
	_, err := QueryOneTx[string](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, pgx.GenericSqlError))

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestIT_QueryOneTx_WhenNoData_ExpectFailure(t *testing.T) {
	_, tx := newTestTransaction(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	_, err := QueryOneTx[element](context.Background(), tx, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, NoMatchingRows))
}

func TestIT_QueryOneTx_WhenTooManyRows_ExpectFailure(t *testing.T) {
	_, tx := newTestTransaction(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table"
	_, err := QueryOneTx[element](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, TooManyMatchingRows))
}

func TestIT_QueryOneTx_ToStruct(t *testing.T) {
	_, tx := newTestTransaction(t)

	type element struct {
		Id   string
		Name string
	}

	sqlQuery := "SELECT id, name FROM my_table WHERE name = 'test-name'"
	actual, err := QueryOneTx[element](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := element{
		Id:   "0463ed3d-bfc9-4c10-b6ee-c223bbca0fab",
		Name: "test-name",
	}
	assert.Equal(expected, actual)
}

func TestIT_QueryOneTx_ToString(t *testing.T) {
	_, tx := newTestTransaction(t)

	sqlQuery := "SELECT id FROM my_table WHERE name = 'test-name'"
	actual, err := QueryOneTx[string](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab", actual)
}

func TestIT_QueryOneTx_ToUuid(t *testing.T) {
	_, tx := newTestTransaction(t)

	sqlQuery := "SELECT id FROM my_table WHERE name = 'test-name'"
	actual, err := QueryOneTx[uuid.UUID](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := uuid.MustParse("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab")
	assert.Equal(expected, actual)
}

func TestIT_QueryAllTx_UnsupportedConnection(t *testing.T) {
	_, err := QueryAllTx[int](context.Background(), &dummyTransaction{}, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, UnsupportedOperation))
}

func TestIT_QueryAllTx_WhenCommitted_ExpectFailure(t *testing.T) {
	_, tx := newTestTransaction(t)
	tx.Close(context.Background())

	_, err := QueryAllTx[int](context.Background(), tx, sampleSqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, AlreadyCommitted))
}

func TestIT_QueryAllTx_WhenConnectionFails_ExpectFailure(t *testing.T) {
	_, tx := newTestTransaction(t)

	sqlQuery := "SELECT name FROM my_tables"
	_, err := QueryAllTx[string](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.NotNil(err)
	assert.True(errors.IsErrorWithCode(err, pgx.GenericSqlError))

	cause := errors.Unwrap(err)
	assert.NotNil(cause)
}

func TestIT_QueryAllTx_NoData(t *testing.T) {
	_, tx := newTestTransaction(t)

	type element struct {
		Name string
	}

	sqlQuery := "SELECT name FROM my_table WHERE name = $1"
	out, err := QueryAllTx[element](context.Background(), tx, sqlQuery, "does-not-exist")

	assert := assert.New(t)
	assert.Nil(err)
	assert.Empty(out)
}

func TestIT_QueryAllTx_ToStruct(t *testing.T) {
	_, tx := newTestTransaction(t)

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
	actual, err := QueryAllTx[element](context.Background(), tx, sqlQuery)

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

func TestIT_QueryAllTx_ToString(t *testing.T) {
	_, tx := newTestTransaction(t)

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
	actual, err := QueryAllTx[string](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := []string{
		"0463ed3d-bfc9-4c10-b6ee-c223bbca0fab",
		"09dd5fc3-0732-4017-81e0-ffee3211d2b9",
	}
	assert.Equal(expected, actual)
}

func TestIT_QueryAllTx_ToUuid(t *testing.T) {
	_, tx := newTestTransaction(t)

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
	actual, err := QueryAllTx[uuid.UUID](context.Background(), tx, sqlQuery)

	assert := assert.New(t)
	assert.Nil(err)
	expected := []uuid.UUID{
		uuid.MustParse("0463ed3d-bfc9-4c10-b6ee-c223bbca0fab"),
		uuid.MustParse("09dd5fc3-0732-4017-81e0-ffee3211d2b9"),
	}
	assert.Equal(expected, actual)
}
