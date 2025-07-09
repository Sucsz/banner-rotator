package postgres

import (
	"context"
	"fmt"

	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/log"
	"github.com/jackc/pgx/v5"
)

// Init подключается к базе данных PostgreSQL и возвращает соединение.
func Init(cfg config.PostgresConfig) (*pgx.Conn, error) {
	logger := log.WithComponent("postgres")

	// Контекст с таймаутом из конфига
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to open PostgreSQL connection.")
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	logger.Info().Msg("Connected to PostgreSQL.")
	return conn, nil
}

// Close закрывает соединение с PostgreSQL.
func Close(ctx context.Context, conn *pgx.Conn) {
	logger := log.WithComponent("postgres")
	if conn == nil {
		return
	}

	if err := conn.Close(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to close PostgreSQL connection.")
	} else {
		logger.Info().Msg("PostgreSQL connection closed.")
	}
}
