package server

import "time"

type Config struct {
	Port            uint16
	ShutdownTimeout time.Duration
}
