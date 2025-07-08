package main

import (
	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/log"
)

func main() {
	// Дефолтный лог
	log.Init("info")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithComponent("main").Fatal().Err(err).
			Msg("Failed to load config")
	}

	// Лог из конфига
	log.Init(cfg.LogLevel)
	log.WithComponent("main").Info().
		Msgf("Service starting on port %s (log level=%s)", cfg.HTTPPort, cfg.LogLevel)
	// далее: connect to DB, Kafka, запустить HTTP-сервер и т.д.
}
