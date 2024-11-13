package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/server"
)

func main() {
	config := server.Config{
		Port:            1234,
		ShutdownTimeout: 2 * time.Second,
	}

	s := server.New(config)
	err := s.Start(context.Background())
	fmt.Printf("err: %v\n", err)
}
