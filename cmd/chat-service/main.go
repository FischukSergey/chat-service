package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	keycloakclient "github.com/FischukSergey/chat-service/internal/clients/keycloak"
	"github.com/FischukSergey/chat-service/internal/config"
	"github.com/FischukSergey/chat-service/internal/logger"
	messagesrepo "github.com/FischukSergey/chat-service/internal/repositories/messages"
	clientv1 "github.com/FischukSergey/chat-service/internal/server-client/v1"
	serverdebug "github.com/FischukSergey/chat-service/internal/server-debug"
	"github.com/FischukSergey/chat-service/internal/store"
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

	// Фикс: В main предположительно должен появиться следующий код:
	// Фикс: 1) Создание клиента к PSQL (не забудьте его "закрыть")
	storeClient, err := store.NewPSQLClient(
		store.NewPSQLOptions(
			cfg.Stores.PSQL.Address,
			cfg.Stores.PSQL.Username,
			cfg.Stores.PSQL.Password,
			cfg.Stores.PSQL.Database,
			cfg.Stores.PSQL.Debug,
		),
	)
	if err != nil {
		return fmt.Errorf("init store client: %v", err)
	}
	defer func() {
		if err := storeClient.Close(); err != nil {
			zap.L().Error("close store client", zap.Error(err))
		}
	}()

	// Фикс: 2) Миграция: https://entgo.io/docs/migrate/#auto-migration
	if err := storeClient.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("create schema: %v", err)
	}

	// Фикс: 3) Инициализация store.Database
	database := store.NewDatabase(storeClient)

	// Фикс: 4) Инициализация репозитория сообщений
	messagesRepository, err := messagesrepo.New(messagesrepo.NewOptions(database))
	if err != nil {
		return fmt.Errorf("create messages repository: %v", err)
	}

	// Загружаем Swagger спецификацию
	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("loading swagger spec: %w", err)
	}
	// Очищаем серверы из спецификации для избежания конфликтов
	swagger.Servers = nil

	// Инициализируем Keycloak клиент
	keycloakClient, err := initKeycloakClient(cfg.Clients.Keycloak, cfg.Global.Env)
	if err != nil {
		return fmt.Errorf("init keycloak client: %v", err)
	}

	// init debug server
	srvDebug, err := serverdebug.New(serverdebug.NewOptions(cfg.Servers.Debug.Addr))
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	// init server client
	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		swagger,
		keycloakClient,
		messagesRepository,
	)
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

// initKeycloakClient инициализирует клиент для Keycloak.
func initKeycloakClient(cfg config.KeycloakConfig, env string) (*keycloakclient.Client, error) {
	lg := zap.L().Named("keycloak-client")

	// Проверяем, если DebugMode включен и окружение prod, предупреждаем об этом
	if cfg.DebugMode && env == "prod" {
		lg.Warn("keycloak client is running in debug mode in production environment")
	}

	// Создаем клиент Keycloak
	client, err := keycloakclient.New(keycloakclient.NewOptions(
		cfg.BasePath,
		cfg.Realm,
		cfg.ClientID,
		cfg.ClientSecret,
		keycloakclient.WithDebugMode(cfg.DebugMode),
	))
	if err != nil {
		return nil, fmt.Errorf("create keycloak client: %v", err)
	}

	return client, nil
}
