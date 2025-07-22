// Package migrator содержит мигратор для применения автоматических миграций с помощью гусь
package migrator

import (
	"database/sql"
	"fmt"

	// Инициализируем драйвер PostgreSQL для работы миграций через goose.
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/log"
)

// Run запускает все миграции "up" из каталога migrations.
func Run(cfg *config.Config) error {
	logger := log.WithComponent("migrate")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("migrator.Run: sql.Open: %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error().Err(err).Msg("migrator: failed to close DB")
		}
	}(db)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrator.Run: SetDialect: %w", err)
	}
	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		return fmt.Errorf("migrator.Run: Up: %w", err)
	}
	return nil
}
