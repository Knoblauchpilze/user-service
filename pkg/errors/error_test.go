package errors

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errSomeError = fmt.Errorf("some error")

const someCode = ErrorCode(26)

func TestUnit_Error_New(t *testing.T) {
	assert := assert.New(t)

	err := New("haha")

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("haha", impl.Message)
	assert.Nil(impl.Cause)
	assert.Equal(GenericErrorCode, impl.Value)
}

func TestUnit_Error_NewCode(t *testing.T) {
	assert := assert.New(t)

	err := NewCode(someCode)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("An unexpected error occurred", impl.Message)
	assert.Nil(impl.Cause)
	assert.Equal(someCode, impl.Value)
}

func TestUnit_Error_NewCode_WhenGenericCode_ExpectMessage(t *testing.T) {
	testCases := []struct {
		code            ErrorCode
		expectedMessage string
	}{
		{code: NotImplementedCode, expectedMessage: "Not implemented"},
		{code: GenericErrorCode, expectedMessage: "An unexpected error occurred"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedMessage, func(t *testing.T) {
			err := NewCode(testCase.code)
			impl, ok := err.(errorImpl)
			assert.True(t, ok)
			assert.Equal(t, testCase.expectedMessage, impl.Message)
		})
	}
}

func TestUnit_Error_NewNotImplemented(t *testing.T) {
	assert := assert.New(t)

	err := NotImplemented()

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("Not implemented", impl.Message)
	assert.Nil(impl.Cause)
	assert.Equal(NotImplementedCode, impl.Value)
}

func TestUnit_Error_NewCodeWithDetails(t *testing.T) {
	assert := assert.New(t)

	err := NewCodeWithDetails(someCode, "message")

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("message", impl.Message)
	assert.Nil(impl.Cause)
	assert.Equal(someCode, impl.Value)
}

func TestUnit_Error_Newf(t *testing.T) {
	assert := assert.New(t)

	err := Newf("haha %d", 22)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("haha 22", impl.Message)
	assert.Nil(impl.Cause)
}

func TestUnit_Error_Wrap(t *testing.T) {
	assert := assert.New(t)

	err := Wrap(errSomeError, "context")

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("context", impl.Message)
	assert.Equal(errSomeError, impl.Cause)
}

func TestUnit_Error_WrapCode(t *testing.T) {
	assert := assert.New(t)

	err := WrapCode(errSomeError, someCode)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("An unexpected error occurred", impl.Message)
	assert.Equal(errSomeError, impl.Cause)
	assert.Equal(someCode, impl.Value)
}

func TestUnit_Error_Wrapf(t *testing.T) {
	assert := assert.New(t)

	err := Wrapf(errSomeError, "context %d", -44)

	impl, ok := err.(errorImpl)
	assert.True(ok)
	assert.Equal("context -44", impl.Message)
	assert.Equal(errSomeError, impl.Cause)
}

func TestUnit_Error_Unwrap(t *testing.T) {
	assert := assert.New(t)

	err := Unwrap(nil)
	assert.Nil(err)

	err = Unwrap(errSomeError)
	assert.Nil(err)

	err = New("haha")
	cause := Unwrap(err)
	assert.Nil(cause)

	err = Wrap(errSomeError, "haha")
	cause = Unwrap(err)
	assert.Equal(errSomeError, cause)

	causeOfCause := Unwrap(cause)
	assert.Nil(causeOfCause)
}

func TestUnit_Error_Error(t *testing.T) {
	assert := assert.New(t)

	err := Wrapf(errSomeError, "context %d", -44)

	expected := "(1) context -44 (cause: some error)"
	assert.Equal(expected, err.Error())

	err = WrapCode(errSomeError, someCode)

	expected = "(26) An unexpected error occurred (cause: some error)"
	assert.Equal(expected, err.Error())
}

func TestUnit_Error_Code(t *testing.T) {
	assert := assert.New(t)

	err := NewCode(someCode)

	impl, ok := err.(ErrorWithCode)
	assert.True(ok)
	assert.Equal(someCode, impl.Code())
}

func TestUnit_Error_MarshalJSON(t *testing.T) {
	assert := assert.New(t)

	err := New("haha")
	out, mErr := json.Marshal(err)

	expected := `
	{
		"Code": 1,
		"Message": "haha"
	}`
	assert.Nil(mErr)
	assert.JSONEq(expected, string(out))

	err = NewCode(someCode)
	out, mErr = json.Marshal(err)

	assert.Nil(mErr)
	expected = `
	{
		"Code": 26,
		"Message": "An unexpected error occurred"
	}`
	assert.JSONEq(expected, string(out))

	err = Wrap(errSomeError, "hihi")
	out, mErr = json.Marshal(err)

	expected = `
	{
		"Code": 1,
		"Message": "hihi",
		"Cause": "some error"
	}`
	assert.Nil(mErr)
	assert.JSONEq(expected, string(out))

	err = Wrap(New("haha"), "hihi")
	out, mErr = json.Marshal(err)

	expected = `
	{
		"Code": 1,
		"Message": "hihi",
		"Cause": {
			"Code": 1,
			"Message": "haha"
		}
	}`
	assert.Nil(mErr)
	assert.JSONEq(expected, string(out))
}

func TestUnit_Error_IsErrorWithCode(t *testing.T) {
	assert := assert.New(t)

	assert.False(IsErrorWithCode(nil, someCode))
	assert.False(IsErrorWithCode(errSomeError, someCode))
	assert.True(IsErrorWithCode(NewCode(someCode), someCode))
	assert.False(IsErrorWithCode(NewCode(27), someCode))
	assert.False(IsErrorWithCode(New("haha"), someCode))
}
