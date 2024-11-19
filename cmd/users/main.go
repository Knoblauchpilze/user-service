package main

import (
	"context"
	"os"
	"time"

	"github.com/KnoblauchPilze/user-service/internal/controller"
	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/KnoblauchPilze/user-service/pkg/server"
)

func main() {
	log := logger.New(logger.NewPrettyWriter(os.Stdout))

	dbName := "my-db"
	dbUser := "my-user"
	dbPassword := "my-password"
	dbConfig := postgresql.NewConfigForLocalhost(dbName, dbUser, dbPassword)
	conn, err := db.New(context.Background(), dbConfig)
	if err != nil {
		log.Errorf("Failed to create db connection: %v", err)
		os.Exit(1)
	}

	repos := repositories.Repositories{
		User: repositories.NewUserRepository(conn),
	}

	apiKeyConfig := service.ApiKeyConfig{
		Validity: 5 * time.Minute,
	}
	userService := service.NewUserService(apiKeyConfig, conn, repos)

	config := server.Config{
		Port:            1234,
		ShutdownTimeout: 2 * time.Second,
	}
	s := server.NewWithLogger(config, log)

	for _, route := range controller.UserEndpoints(userService) {
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
