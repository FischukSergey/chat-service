package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/FischukSergey/chat-service/internal/clients/keycloak"
	serverclient "github.com/FischukSergey/chat-service/internal/server-client"
	clientv1 "github.com/FischukSergey/chat-service/internal/server-client/v1"
)

const nameServerClient = "server-client"

func initServerClient( // воспользуйся мной в chat-service/main.go
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	keycloakIntrospector *keycloakclient.Client,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	// Создаем опции для сервера
	options := []serverclient.OptOptionsSetter{}

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
