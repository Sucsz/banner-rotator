package log

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init Конфигурирует Zerolog Global Logger на основе предоставленного уровня.
func Init(level string) {
	// Set timestamp format to RFC3339
	zerolog.TimeFieldFormat = time.RFC3339

	// форматируем JSON-лог в человекочитаемый вид в консоль
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	log.Logger = log.Output(console)

	// Используем переданный уровень логирования(из конфига)
	lvl, err := zerolog.ParseLevel(level)
	// Если передали что-то не то, то оставляем дефолтный уровень инфо
	if err != nil {
		log.Warn().Err(err).Msg("Invalid log level, defaulting to info")
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
}

// WithComponent позволяет создавать контекстные логгеры для каждого модуля для упрощения фильтрации по компонентам
func WithComponent(name string) *zerolog.Logger {
	logger := zerolog.Ctx(context.Background()).With().
		Str("component", name).
		Logger()
	return &logger
}
