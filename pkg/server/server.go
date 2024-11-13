package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/rest"
	"github.com/labstack/echo/v4"
)

type Server interface {
	AddRoute(route rest.Route) error
	Start(ctx context.Context) error
}

type serverImpl struct {
	echo            *echo.Echo
	port            uint16
	shutdownTimeout time.Duration
	router          *echo.Group
}

func New(config Config) Server {
	echo := createEchoServer()

	s := &serverImpl{
		echo:            echo,
		port:            config.Port,
		shutdownTimeout: config.ShutdownTimeout,
		router:          echo.Group(""),
	}

	return s
}

func (s *serverImpl) AddRoute(route rest.Route) error {
	path := route.Path()

	switch route.Method() {
	case http.MethodGet:
		s.router.GET(path, route.Handler())
	case http.MethodPost:
		s.router.POST(path, route.Handler())
	case http.MethodDelete:
		s.router.DELETE(path, route.Handler())
	case http.MethodPatch:
		s.router.PATCH(path, route.Handler())
	default:
		return errors.NewCode(UnsupportedMethod)
	}

	return nil
}

func (s *serverImpl) Start(ctx context.Context) error {
	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	notifyCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	var runError error

	go func() {
		address := fmt.Sprintf(":%d", s.port)
		err := s.echo.Start(address)

		if err != nil && err != http.ErrServerClosed {
			runError = err
		}
	}()

	const reasonableWaitTimeToInitializeServer = 50 * time.Millisecond
	time.Sleep(reasonableWaitTimeToInitializeServer)

	<-notifyCtx.Done()

	err := s.shutdown()
	if err != nil {
		return err
	}
	return runError
}

func (s *serverImpl) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.echo.Shutdown(ctx)
}

func createEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return e
}
