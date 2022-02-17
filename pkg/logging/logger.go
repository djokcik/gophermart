package logging

import (
	"github.com/rs/zerolog"
	"os"
)

// NewLogger creates a new customizable logger.
func NewLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Output(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	return &logger
}
