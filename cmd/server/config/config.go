package config

import (
	"github.com/KnoblauchPilze/user-service/cmd/server/server"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   server.Config
	Database db.Config
}

func Load() (Configuration, error) {
	// https://github.com/spf13/viper#reading-config-files
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	viper.SetConfigName("server-dev")
	if err := viper.ReadInConfig(); err != nil {
		return defaultConf(), err
	}

	var out Configuration
	if err := viper.Unmarshal(&out); err != nil {
		return defaultConf(), err
	}

	return out, nil
}

func defaultConf() Configuration {
	return Configuration{
		Server: server.Config{
			Endpoint: "/v1/users/",
			Port:     uint16(60000),
		},
		Database: db.Config{
			Host:                "localhost",
			Port:                5432,
			Name:                "database",
			User:                "user",
			Password:            "password",
			ConnectionsPoolSize: 1,
		},
	}
}
