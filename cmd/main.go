package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/log"
	"github.com/Sucsz/banner-rotator/pkg/postgres"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализируем логгер
	log.Init(cfg.LogLevel)
	logger := log.WithComponent("main")
	logger.Info().
		Msgf("Service starting on port %s (log level = %s).", cfg.HTTPPort, cfg.LogLevel)

	// Подключаемся к PostgreSQL
	conn, err := postgres.Init(cfg.Postgres)
	if err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to initialize PostgreSQL.")
	}
	// Гарантированно закроем соединение при выходе
	defer postgres.Close(context.Background(), conn)

	// TODO: здесь будет запуск HTTP-сервера и остальных компонентов
}
