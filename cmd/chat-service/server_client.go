package main

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	keycloakclient "github.com/FischukSergey/chat-service/internal/clients/keycloak"
	messagesrepo "github.com/FischukSergey/chat-service/internal/repositories/messages"
	serverclient "github.com/FischukSergey/chat-service/internal/server-client"
	clientv1 "github.com/FischukSergey/chat-service/internal/server-client/v1"
	gethistory "github.com/FischukSergey/chat-service/internal/usecases/client/get-history"
)

const nameServerClient = "server-client"

func initServerClient( // воспользуйся мной в chat-service/main.go
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	keycloakIntrospector *keycloakclient.Client,
	messagesRepository *messagesrepo.Repo,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	// Фикс: 1) Создание getHistoryUseCase (и его проброс в хендлеры)
	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(messagesRepository))
	if err != nil {
		return nil, fmt.Errorf("create get history use case: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	// Создаем опции для сервера
	options := []serverclient.OptOptionsSetter{
		serverclient.WithEchoHTTPErrorHandler(initHTTPErrorHandler(lg)),
	}

	// Добавляем опцию для Keycloak, если клиент определен
	if keycloakIntrospector != nil {
		options = append(options, serverclient.WithKeycloakIntrospector(keycloakIntrospector))
	}

	// Создаем сервер
	srv, err := serverclient.New(serverclient.NewOptions(
		lg,
		addr,
		allowOrigins,
		v1Swagger,
		v1Handlers,
		options...,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}

// Фикс: 2) Инициализация httpErrorHandler (и его проброс в сервер)
// Фикс: 3) "server-client" logger должен пронизывать все компоненты сервера
func initHTTPErrorHandler(lg *zap.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		lg.Error("http error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}
