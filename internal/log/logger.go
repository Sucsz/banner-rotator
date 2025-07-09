package log

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var baseLogger zerolog.Logger

// Init конфигурирует глобальный zerolog.Logger по уровню level.
func Init(level string) {
	zerolog.TimeFieldFormat = time.RFC3339

	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	lvl, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		// здесь можно использовать log.Logger, хотя он ещё не переназначен – это мелочь
		log.Warn().Err(err).Msg("Invalid log level, defaulting to info.")
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	baseLogger = zerolog.New(console).With().Timestamp().Logger()
	log.Logger = baseLogger
}

// WithComponent создаёт новый логгер с полем component и возвращает указатель на него,
// чтобы можно было вызывать методы Fatal(), Error(), Info() и т.д.
func WithComponent(name string) *zerolog.Logger {
	l := baseLogger.With().Str("component", name).Logger()
	return &l
}
