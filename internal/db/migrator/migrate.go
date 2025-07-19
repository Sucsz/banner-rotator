package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/Sucsz/banner-rotator/config"
)

// Run запускает все миграции "up" из каталога migrations
func Run(cfg *config.Config) error {
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
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrator.Run: SetDialect: %w", err)
	}
	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		return fmt.Errorf("migrator.Run: Up: %w", err)
	}
	return nil
}
