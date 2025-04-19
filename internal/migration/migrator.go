package migration

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Changed to postgres
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(db *pgxpool.Pool) error {
	log.Println("📦 Running database migrations...")

	// Get connection string from pool config
	connConfig := db.Config().ConnConfig
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		connConfig.User,
		connConfig.Password,
		connConfig.Host,
		connConfig.Port,
		connConfig.Database,
		connConfig.RuntimeParams["sslmode"],
	)

	m, err := migrate.New(
		"file://migrations",
		connString,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Database migration completed.")
	return nil
}
