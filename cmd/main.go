package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nyzhehorodov/apicompanies/api"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/httpserver"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/log"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/log/zap"

	zapoptions "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log.SetLogger(zap.New(func(opts *zap.Options) {
		opts.Development = true
		logLvl := zapoptions.NewAtomicLevelAt(zapcore.DebugLevel)
		opts.Level = &logLvl
	}))

	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())

	app, err := api.New(httpserver.New(), log.Logger.WithName("api"))
	check(err)

	go app.Run(ctx, api.Config{Port: 8080})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)

		cancel()
	}()

	check(<-errCh)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
