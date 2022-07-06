package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nyzhehorodov/apicompanies/api"
	"github.com/nyzhehorodov/apicompanies/pkg/config"
	"github.com/nyzhehorodov/apicompanies/pkg/di"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/httpserver"

	"github.com/spf13/viper"
)

func main() {
	conf, err := initConfig()
	check("init config", err)

	c := di.New("server", conf)
	if conf.Database.Migration.Enabled {
		err = migrateDB(c)
		check("migrate db", err)

	}

	app, err := initAPI(c)
	check("init api", err)

	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)

		cancel()
	}()

	go serveAPI(ctx, conf, app)

	check("got signal", <-errCh)
}

func initConfig() (config.Config, error) {
	viper.SetConfigFile("./config.yaml")
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	var conf config.Config

	err := viper.ReadInConfig()
	if err != nil {
		return conf, fmt.Errorf("read config file, %w", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, fmt.Errorf("decode config, %w", err)
	}

	return conf, nil
}

func migrateDB(c *di.Container) error {
	migrator, err := c.NewMigrator()
	if err != nil {
		return fmt.Errorf("new migrator: %w", err)
	}

	if err = migrator.Migrate(context.Background()); err != nil {
		return fmt.Errorf("exec migrations: %w", err)
	}

	if err := migrator.Close(context.Background()); err != nil {
		return fmt.Errorf("close conn: %w", err)
	}

	return nil
}

func initAPI(c *di.Container) (*api.API, error) {
	companyService, err := c.CompanyService()
	if err != nil {
		return nil, fmt.Errorf("new compane repo: %w", err)
	}

	a := &api.API{
		Server:         httpserver.New(),
		Logger:         c.Logger().WithName("apicompany"),
		CompanyService: companyService,
	}
	a.Init()

	return a, nil
}

func serveAPI(ctx context.Context, conf config.Config, api *api.API) {
	go func() {
		<-ctx.Done()

		if err := api.Server.Shutdown(context.Background()); err != nil {
			api.Logger.Error(err, "server shutdown")
		}
	}()

	listenAddr := fmt.Sprintf(":%d", conf.Server.Port)

	api.Logger.Info("server listen", "addr", listenAddr)

	var err error
	if conf.Server.TLSCert != "" && conf.Server.TLSKey != "" {
		err = api.Server.ListenTLS(listenAddr, conf.Server.TLSCert, conf.Server.TLSKey)
	} else {
		err = api.Server.Listen(listenAddr)
	}
	if !errors.Is(err, http.ErrServerClosed) {
		api.Logger.Error(err, "server listen")
	}
}

func check(msg string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", msg, err))
	}
}
