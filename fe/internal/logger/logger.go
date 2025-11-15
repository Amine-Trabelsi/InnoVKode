package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	log := zerolog.New(output).Level(lvl).With().Timestamp().Logger()
	return log
}
