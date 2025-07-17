package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/Sucsz/banner-rotator/config"
	"github.com/Sucsz/banner-rotator/internal/log"
)

// Run запускает все миграции "up" из каталога migrationsDir.
func Run(cfg *config.Config) error {
	logger := log.WithComponent("migrator")

	// Собираем DSN из конфига (в том же формате, что и в postgres.Init)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	logger.Info().Msgf("Connecting for migrations: %s", cfg.Postgres.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.WithComponent("postgres").Error().Err(err).Msg("Failed to close migrator DB")
		}
	}()

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	logger.Info().Msg("Migrations applied successfully")
	return nil
}
