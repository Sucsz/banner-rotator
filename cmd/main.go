package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/db/migrator"
	"github.com/Sucsz/banner-rotator/internal/log"
	"github.com/Sucsz/banner-rotator/pkg/postgres"
)

func main() {
	// 1) Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2) Инициализируем логгер
	log.Init(cfg.LogLevel)
	logger := log.WithComponent("main")
	logger.Info().
		Msgf("Service starting on port %s (log level = %s).", cfg.HTTPPort, cfg.LogLevel)

	// 3) Прогон миграций
	if err := migrator.Run(cfg); err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to run database migrations.")
	}

	// 4) Подключаемся к PostgreSQL
	conn, err := postgres.Init(cfg.Postgres)
	if err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to initialize PostgreSQL.")
	}
	// Гарантированно закроем соединение при выходе
	defer postgres.Close(context.Background(), conn)

	// TODO: здесь будет запуск HTTP‑сервера, Kafka‑producer и т.д.
}
