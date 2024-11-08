package postgresql

import (
	"time"
)

type Config struct {
	Host           string
	Port           uint16
	Database       string
	User           string
	Password       string
	ConnectTimeout time.Duration
}

const defaultConnectTimeout = 5 * time.Second

func (c Config) ToConnectionString() string {
	return generateConnectionString(c)
}

func NewConfigForLocalhost(database string, user string, password string) Config {
	return Config{
		Host:           "localhost",
		Port:           5432,
		Database:       database,
		User:           user,
		Password:       password,
		ConnectTimeout: defaultConnectTimeout,
	}
}
