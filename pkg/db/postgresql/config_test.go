package postgresql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Config_ExpectCorrectConnectionString(t *testing.T) {
	config := Config{
		Host:           "my-host",
		Port:           8976,
		Database:       "my-database",
		User:           "super-user",
		Password:       "strong-password?",
		ConnectTimeout: 2 * time.Hour,
	}

	actual := config.ToConnectionString()

	assert := assert.New(t)
	expected := "postgresql://super-user:strong-password%3F@my-host:8976/my-database?connect_timeout=7200"
	assert.Equal(expected, actual)
}
