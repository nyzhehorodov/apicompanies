package migration

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"

	"github.com/nyzhehorodov/apicompanies/pkg/lib/log"
)

const defaultVersionTable = "schema_version"

type Config struct {
	MigrationsPath string
	VersionTable   string
}

type Migrator interface {
	Migrate(ctx context.Context) error
	MigrateTo(ctx context.Context, targetVersion int32) error
	GetCurrentVersion(ctx context.Context) (v int32, err error)
	Close(ctx context.Context) error
}

func New(conf Config, conn *pgx.Conn, logger log.Interface) (Migrator, error) {
	if conf.VersionTable == "" {
		conf.VersionTable = defaultVersionTable
	}

	migrator, err := migrate.NewMigrator(context.Background(), conn, conf.VersionTable)
	if err != nil {
		return nil, fmt.Errorf("initializing migrator: %w", err)
	}

	err = migrator.LoadMigrations(conf.MigrationsPath)
	if err != nil {
		return nil, fmt.Errorf("loading migration: %w", err)
	}

	migrator.OnStart = func(sequence int32, name, direction, sql string) {
		logger.Info(fmt.Sprintf(
			"executing %s %s\n%s",
			name,
			direction,
			sql,
		))
	}

	return &impl{
		migrator: migrator,
		conn:     conn,
	}, nil
}

type impl struct {
	migrator *migrate.Migrator
	conn     *pgx.Conn
}

func (i impl) Migrate(ctx context.Context) error {
	return i.migrator.Migrate(ctx)
}

func (i impl) MigrateTo(ctx context.Context, targetVersion int32) error {
	return i.migrator.MigrateTo(ctx, targetVersion)
}

func (i impl) GetCurrentVersion(ctx context.Context) (v int32, err error) {
	return i.migrator.GetCurrentVersion(ctx)
}

func (i impl) Close(ctx context.Context) error {
	return i.conn.Close(ctx)
}
