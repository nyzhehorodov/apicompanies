package di

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	zapoptions "go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nyzhehorodov/apicompanies/pkg/app/company"
	"github.com/nyzhehorodov/apicompanies/pkg/config"
	dcompany "github.com/nyzhehorodov/apicompanies/pkg/domain/company"
	"github.com/nyzhehorodov/apicompanies/pkg/infra/db"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/log"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/log/zap"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/migration"
)

// Container is a dependency injection container to be used in the main packages.
// All common dependency initialization should go here.
type Container struct {
	name           string
	conf           config.Config
	log            log.Interface
	connPool       *pgxpool.Pool
	companyRepo    dcompany.Repository
	companyService company.Service
}

func New(name string, conf config.Config) *Container {
	return &Container{
		name: name,
		conf: conf,
	}
}

func (c *Container) Logger() log.Interface {
	if c.log.Enabled() {
		return c.log
	}

	log.SetLogger(zap.New(func(opts *zap.Options) {
		opts.Development = c.conf.Log.Development
		logLvl := zapoptions.NewAtomicLevelAt(-zapcore.Level(c.conf.Log.Verbosity))
		opts.Level = &logLvl
	}))

	c.log = log.Logger

	c.log.Info("logger", "verbosity", c.conf.Log.Verbosity, "development", c.conf.Log.Development)
	return c.log
}

func (c *Container) ConnPool() (*pgxpool.Pool, error) {
	if c.connPool != nil {
		return c.connPool, nil
	}

	conf, err := pgxpool.ParseConfig(c.conf.Database.URI)
	if err != nil {
		return nil, fmt.Errorf("pgx parse config: %w", err)
	}

	if c.conf.Database.MinConns > 0 {
		conf.MinConns = c.conf.Database.MinConns
	}
	if c.conf.Database.MaxConns > 0 {
		conf.MaxConns = c.conf.Database.MaxConns
	}

	connPool, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("new pgx conn pool: %w", err)
	}

	c.connPool = connPool

	return c.connPool, nil
}

func (c *Container) NewMigrator() (migration.Migrator, error) {
	conn, err := pgx.Connect(context.Background(), c.conf.Database.URI)
	if err != nil {
		return nil, fmt.Errorf("new pgx conn: %w", err)
	}

	migrator, err := migration.New(
		migration.Config{
			MigrationsPath: c.conf.Database.Migration.Path,
			VersionTable:   c.conf.Database.Migration.VersionTable,
		},
		conn,
		c.Logger().WithName("migration"),
	)
	if err != nil {
		return nil, fmt.Errorf("new migrator: %w", err)
	}

	return migrator, nil
}

func (c *Container) CompanyService() (company.Service, error) {
	if c.companyService != nil {
		return c.companyService, nil
	}

	companyRepo, err := c.CompanyRepo()
	if err != nil {
		return nil, err
	}

	c.companyService = company.NewService(companyRepo)

	return c.companyService, nil
}

func (c *Container) CompanyRepo() (dcompany.Repository, error) {
	if c.companyRepo != nil {
		return c.companyRepo, nil
	}

	conn, err := c.ConnPool()
	if err != nil {
		return nil, err
	}

	c.companyRepo = db.NewCompanyPostgresRepository(conn)

	return c.companyRepo, nil
}
