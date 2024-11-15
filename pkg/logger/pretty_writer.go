package logger

import (
	"io"
	"time"

	"github.com/rs/zerolog"
)

func NewPrettyWriter(out io.Writer) io.Writer {
	return zerolog.ConsoleWriter{
		Out:          out,
		TimeFormat:   time.DateTime,
		TimeLocation: time.UTC,
	}
}
