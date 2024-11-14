package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Format(t *testing.T) {
	type testCase struct {
		color        Color
		expectedText string
	}

	testCases := map[string]testCase{
		"blue": {

			color:        Blue,
			expectedText: "\033[1;34mhello\033[0m",
		},
		"cyan": {
			color:        Cyan,
			expectedText: "\033[1;36mhello\033[0m",
		},
		"gray": {
			color:        Gray,
			expectedText: "\033[1;90mhello\033[0m",
		},
		"green": {
			color:        Green,
			expectedText: "\033[1;32mhello\033[0m",
		},
		"magenta": {
			color:        Magenta,
			expectedText: "\033[1;35mhello\033[0m",
		},
		"yellow": {
			color:        Yellow,
			expectedText: "\033[1;33mhello\033[0m",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := FormatWithColor("hello", testCase.color)

			assert := assert.New(t)
			assert.Equal(testCase.expectedText, actual)
		})
	}
}

func TestFormat_Red(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Red)
	assert.Equal("\033[1;31mhello\033[0m", actual)
}

func TestFormat_Greeb(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Green)
	assert.Equal("\033[1;32mhello\033[0m", actual)
}

func TestFormat_Yellow(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Yellow)
	assert.Equal("\033[1;33mhello\033[0m", actual)
}

func TestFormat_Blue(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Blue)
	assert.Equal("\033[1;34mhello\033[0m", actual)
}

func TestFormat_Magenta(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Magenta)
	assert.Equal("\033[1;35mhello\033[0m", actual)
}

func TestFormat_Cyan(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Cyan)
	assert.Equal("\033[1;36mhello\033[0m", actual)
}

func TestFormat_Gray(t *testing.T) {
	assert := assert.New(t)

	actual := FormatWithColor("hello", Gray)
	assert.Equal("\033[1;90mhello\033[0m", actual)
}
