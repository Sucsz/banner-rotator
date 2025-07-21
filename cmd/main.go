package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/api"
	"github.com/Sucsz/banner-rotator/internal/db/dao"
	"github.com/Sucsz/banner-rotator/internal/db/migrator"
	"github.com/Sucsz/banner-rotator/internal/kafka"
	"github.com/Sucsz/banner-rotator/internal/log"
	"github.com/Sucsz/banner-rotator/internal/service/bandit"
	"github.com/Sucsz/banner-rotator/pkg/postgres"
)

func main() {
	// 1) Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2) Инициализируем глобальный логгер
	log.Init(cfg.LogLevel)
	logger := log.WithComponent("main")
	logger.Info().
		Msgf("Service starting on port %s (log level = %s).", cfg.HTTPPort, cfg.LogLevel)

	// 3) Прогон миграций
	if err := migrator.Run(cfg); err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to run database migrations.")
	}
	logger.Info().Msg("Migrations applied successfully.")

	// 4) Подключаемся к PostgreSQL
	conn, err := postgres.Init(cfg.Postgres)
	if err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to initialize PostgreSQL.")
	}
	defer func() {
		if err := postgres.Close(context.Background(), conn); err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to close PostgreSQL connection.")
		}
	}()

	// 5) Проверяем доступность Kafka‑брокера
	if err := kafka.CheckConnection(cfg.Kafka.Brokers, 5*time.Second); err != nil {
		logger.Fatal().
			Err(err).
			Msg("Kafka broker is not reachable.")
	}
	logger.Info().
		Strs("brokers", cfg.Kafka.Brokers).
		Str("topic", cfg.Kafka.Topic).
		Msg("Kafka broker connection successful.")

	// 6) Инициализируем Kafka‑продюсера
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	logger.Info().
		Strs("brokers", cfg.Kafka.Brokers).
		Str("topic", cfg.Kafka.Topic).
		Msg("Kafka producer initialized.")
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to close Kafka producer.")
		}
	}()

	// 7) Инициализируем DAO-слой
	statDAO := dao.NewStatDAO(conn)
	slotDAO := dao.NewBannerSlotDAO(conn)

	// 8) Создаём ε‑greedy селектор
	selector := bandit.NewBandit(
		bandit.Config{Epsilon: cfg.Epsilon},
		statDAO, slotDAO,
	)

	// 9) Собираем API и роутер
	apiHandler := api.NewAPI(selector, producer, slotDAO)
	router := api.NewRouter(apiHandler)

	// 10) Запускаем HTTP-сервер
	logger.Info().
		Str("addr", ":"+cfg.HTTPPort).
		Msg("Starting HTTP server.")

	if err := http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to start HTTP server.")
	}
}
