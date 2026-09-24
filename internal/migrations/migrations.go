package migrations

import (
	"embed"
	"errors"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var FS embed.FS

func Run(dsn string) error {
	src, err := iofs.New(FS, "sql")
	if err != nil {
		return fmt.Errorf("migrations: source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		return fmt.Errorf("migrations: init: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("migrations: up: %w", err)
	}
	return nil

}
