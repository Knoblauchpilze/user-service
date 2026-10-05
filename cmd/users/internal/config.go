package internal

import (
	"time"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db/postgresql"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/Knoblauchpilze/user-service/internal/service"
)

type Configuration struct {
	Server   server.Config
	Port     uint16
	Database postgresql.Config
	ApiKey   service.ApiKeyConfig
}

func DefaultConfig() Configuration {
	const defaultDatabaseName = "db_user_service"
	const defaultDatabaseUser = "user_service_manager"

	return Configuration{
		Server: server.Config{
			BasePath:        "/v1/users",
			ShutdownTimeout: 5 * time.Second,
		},
		Port: uint16(80),
		Database: postgresql.NewConfigForDockerContainer(
			defaultDatabaseName,
			defaultDatabaseUser,
			"comes-from-the-environment",
		),
		ApiKey: service.ApiKeyConfig{
			Validity: time.Duration(3 * time.Hour),
		},
	}
}
