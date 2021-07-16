package logger

import (
	"os"
	"sync"

	"github.com/jarsida/OzonIMDG_Case/config"
	"github.com/rs/zerolog"
)

// Logger is a logger :)
type Logger struct {
	*zerolog.Logger
}

var (
	logger Logger
	once   sync.Once
)

// Чтение конфигурации из среды. Единожды.
func Get() *Logger {
	once.Do(func() {
		zeroLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		cfg := config.Get()
		// Set proper loglevel based on config
		switch cfg.LogLevel {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn", "warning":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "err", "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel) // log info and above by default
		}
		logger = Logger{&zeroLogger}
	})
	return &logger
}
