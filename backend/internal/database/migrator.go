package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"strconv"

	"github.com/jackc/pgx/v5"
	tern "github.com/jackc/tern/v2/migrate"
	"github.com/rs/zerolog"
	"guthub.com/imritik7303/boiler-plate-backend/internal/config"
)

//go:embed migrations/*.sql
var migrations embed.FS

//tells Go (via the embed package) to include all .sql files inside the migrations/ directory at compile time.

//embed.FS is a virtual filesystem containing those files.

//This means when you ship the binary, it already contains the migration scripts — no need to copy them separately.

func Migrate(ctx context.Context, logger *zerolog.Logger, cfg *config.Config) error {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	// URL-encode the password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return fmt.Errorf("constructing database migrator: %w", err)
	}
	//Creates a Migrator object from tern.

	//schema_version is the tracking table in Postgres that records which migrations have been applied.

	//If the table doesn’t exist, tern creates it automatically.
	subtree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("retrieving database migrations subtree: %w", err)
	}

	//fs.Sub(migrations, "migrations") → gets a subdirectory (migrations/) from the embedded FS.

	if err := m.LoadMigrations(subtree); err != nil {
		return fmt.Errorf("loading database migrations: %w", err)
	}

	//m.LoadMigrations(subtree) → loads .sql migration files into the migrator.

	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("retreiving current database migration version")
	}
	//Asks tern what the current schema version is (based on the schema_version table).

	if err := m.Migrate(ctx); err != nil {
		return err
	}
	//Applies all unapplied migrations in order, inside transactions.

	if from == int32(len(m.Migrations)) {
		logger.Info().Msgf("database schema up to date, version %d", len(m.Migrations))
	} else {
		logger.Info().Msgf("migrated database schema, from %d to %d", from, len(m.Migrations))
	}
	return nil
}
