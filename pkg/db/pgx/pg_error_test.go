package pgx

import (
	"fmt"
	"testing"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	jpgx "github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestUnit_AnalyzeAndWrapPgError_Nil(t *testing.T) {
	err := AnalyzeAndWrapPgError(nil)

	assert := assert.New(t)
	assert.Nil(err)
}

func TestUnit_AnalyzeAndWrapPgError_WhenNotAPgError_ExpectUnchanged(t *testing.T) {
	err := fmt.Errorf("some error")

	actual := AnalyzeAndWrapPgError(err)

	assert := assert.New(t)
	assert.Equal(err, actual)
}

func TestUnit_AnalyzeAndWrapPgError(t *testing.T) {
	type testCase struct {
		code          string
		expectedError errors.ErrorCode
	}

	testCases := []testCase{
		{
			code:          "23503",
			expectedError: ForeignKeyValidation,
		},
		{
			code:          "23505",
			expectedError: UniqueConstraintViolation,
		},
		{
			code:          "not-a-code",
			expectedError: GenericSqlError,
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			assert := assert.New(t)

			err := &jpgx.PgError{
				Code: testCase.code,
			}

			actual := AnalyzeAndWrapPgError(err)

			assert.True(errors.IsErrorWithCode(actual, testCase.expectedError), "Actual err: %v", err)
			cause := errors.Unwrap(actual)
			assert.Equal(err, cause)
		})
	}
}
