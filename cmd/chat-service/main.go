package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/FischukSergey/chat-service/internal/config"
	"github.com/FischukSergey/chat-service/internal/logger"
	clientv1 "github.com/FischukSergey/chat-service/internal/server-client/v1"
	serverdebug "github.com/FischukSergey/chat-service/internal/server-debug"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	// logger.Init & logger.Sync
	if err := logger.Init(logger.NewOptions(
		cfg.Log.Level,
		logger.WithDsnSentry(cfg.Sentry.DSN),
		logger.WithEnv(cfg.Global.Env),
	)); err != nil {
		return fmt.Errorf("init logger: %v", err)
	}
	defer logger.Sync()

	// Загружаем Swagger спецификацию
	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("loading swagger spec: %w", err)
	}
	// Очищаем серверы из спецификации для избежания конфликтов
	swagger.Servers = nil

	// init debug server
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}
	// init server client
	srvClient, err := initServerClient(cfg.Servers.Client.Addr, cfg.Servers.Client.AllowOrigins, swagger)
	if err != nil {
		return fmt.Errorf("init server client: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })
	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
