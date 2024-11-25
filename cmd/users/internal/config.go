package internal

import (
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/server"
	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
)

type Configuration struct {
	Server   server.Config
	Database postgresql.Config
	ApiKey   service.ApiKeyConfig
}

func DefaultConfig() Configuration {
	const defaultDatabaseName = "db_user_service"
	const defaultDatabaseUser = "user_service_manager"

	return Configuration{
		Server: server.Config{
			BasePath:        "/v1/users",
			Port:            uint16(80),
			ShutdownTimeout: 5 * time.Second,
		},
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
