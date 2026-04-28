// Package main starts the user-service HTTP server.
//
// @title User Service API
// @version 1.0
// @description HTTP API for managing users, sessions, and API key authentication.
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Api-Key
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/config"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/logger"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/process"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/server"
	docs "github.com/Knoblauchpilze/user-service/api"
	"github.com/Knoblauchpilze/user-service/cmd/users/internal"
	"github.com/Knoblauchpilze/user-service/internal/controller"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/Knoblauchpilze/user-service/pkg/repositories"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func determineConfigName() string {
	if len(os.Args) < 2 {
		return "users-prod.yml"
	}

	return os.Args[1]
}

func main() {
	log := logger.New(os.Stdout)

	conf, err := config.Load(determineConfigName(), internal.DefaultConfig())
	if err != nil {
		log.Error("Failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	conn, err := db.New(context.Background(), conf.Database)
	if err != nil {
		log.Error("Failed to create db connection", slog.Any("error", err))
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
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	for _, route := range controller.HealthCheckEndpoints(conn) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	for _, route := range controller.AuthEndpoints(authService) {
		if err := s.AddRoute(route); err != nil {
			log.Error("Failed to register route", slog.String("route", route.Path()), slog.Any("error", err))
			os.Exit(1)
		}
	}

	docs.SwaggerInfo.BasePath = conf.Server.BasePath
	swaggerUi := rest.NewRawRoute(http.MethodGet, "/swagger/*", echoSwagger.WrapHandler)
	if err := s.AddRoute(swaggerUi); err != nil {
		log.Error("Failed to register route", slog.String("route", swaggerUi.Path()), slog.Any("error", err))
		os.Exit(1)
	}

	wait, err := process.StartWithSignalHandler(context.Background(), s)
	if err != nil {
		log.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}

	err = wait()
	if err != nil {
		log.Error("Error while serving", slog.Any("error", err))
		os.Exit(1)
	}
}
