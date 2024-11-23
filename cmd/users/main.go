package main

import (
	"context"
	"os"

	"github.com/KnoblauchPilze/user-service/cmd/users/internal"
	"github.com/KnoblauchPilze/user-service/internal/controller"
	"github.com/KnoblauchPilze/user-service/internal/service"
	"github.com/KnoblauchPilze/user-service/pkg/config"
	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/logger"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/KnoblauchPilze/user-service/pkg/server"
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
		log.Errorf("Failed to create db connection: %v", err)
		os.Exit(1)
	}

	conn, err := db.New(context.Background(), conf.Database)
	if err != nil {
		log.Errorf("Failed to create db connection: %v", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	repos := repositories.Repositories{
		User: repositories.NewUserRepository(conn),
	}

	userService := service.NewUserService(conf.ApiKey, conn, repos)

	s := server.NewWithLogger(conf.Server, log)

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
