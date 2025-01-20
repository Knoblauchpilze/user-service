package main

import (
	"context"
	"os"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/config"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/logger"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	"github.com/Knoblauchpilze/user-service/cmd/users/internal"
	"github.com/Knoblauchpilze/user-service/internal/controller"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "users-prod.yml"
	}

	return os.Args[1]
}

func main() {
	log := logger.New(logger.NewPrettyWriter(os.Stdout))

	conf, err := config.Load(determineConfigName(), internal.DefaultConfig())
	if err != nil {
		log.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	conn, err := db.New(context.Background(), conf.Database)
	if err != nil {
		log.Errorf("Failed to create db connection: %v", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	repos := repositories.Repositories{
		User:   repositories.NewUserRepository(conn),
		ApiKey: repositories.NewApiKeyRepository(conn),
	}

	userService := service.NewUserService(conf.ApiKey, conn, repos)
	authService := service.NewAuthService(repos)

	s := server.NewWithLogger(conf.Server, log)

	for _, route := range controller.UserEndpoints(userService) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route %v: %v", route.Path(), err)
			os.Exit(1)
		}
	}

	for _, route := range controller.HealthCheckEndpoints(conn) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route %v: %v", route.Path(), err)
			os.Exit(1)
		}
	}

	for _, route := range controller.AuthEndpoints(authService) {
		if err := s.AddRoute(route); err != nil {
			log.Errorf("Failed to register route %v: %v", route.Path(), err)
			os.Exit(1)
		}
	}

	err = s.Start(context.Background())
	if err != nil {
		log.Errorf("Error while serving: %v", err)
		os.Exit(1)
	}
}
