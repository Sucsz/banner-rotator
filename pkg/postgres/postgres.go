package postgres

import (
	"context"
	"fmt"
	"github.com/Sucsz/banner-rotator/config"
	"github.com/jackc/pgx/v5"
)

// Init подключается к базе данных PostgreSQL и возвращает соединение.
func Init(cfg config.PostgresConfig) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres.Init: %w", err)
	}
	return conn, nil
}

// Close закрывает соединение с PostgreSQL.
func Close(ctx context.Context, conn *pgx.Conn) error {
	if conn == nil {
		return nil
	}
	if err := conn.Close(ctx); err != nil {
		return fmt.Errorf("postgres.Close: %w", err)
	}
	return nil
}
