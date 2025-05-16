package serverclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimdlwr "github.com/oapi-codegen/echo-middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	clientv1 "github.com/FischukSergey/chat-service/internal/server-client/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger       *zap.Logger              `option:"mandatory" validate:"required"`
	addr         string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins []string                 `option:"mandatory" validate:"min=1"`
	v1Swagger    *openapi3.T              `option:"mandatory" validate:"required"`
	v1Handlers   clientv1.ServerInterface `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	// валидация опций
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	e := echo.New()
	e.Use(
		// CORS middleware
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.allowOrigins,
			AllowMethods: []string{http.MethodPost, http.MethodOptions},
			AllowHeaders: []string{"X-Request-ID", "Content-Type", "Authorization"},
		}),
	)
	// переделаный авторский вариант ??????????
	// Создаем OpenAPI валидатор с правильными опциями
	validator := oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:  false,
			ExcludeResponseBody: true,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
	})

	// Регистрируем обработчики напрямую на маршрутах без группы v1
	wrapper := &clientv1.ServerInterfaceWrapper{Handler: opts.v1Handlers}
	e.POST("/v1/getHistory", wrapper.PostGetHistory, validator)

	srv := &http.Server{
		Addr:              opts.addr,
		Handler:           e,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return &Server{
		lg:  opts.logger,
		srv: srv,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// Ожидаем завершения контекста и выполняем graceful shutdown
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(shutdownCtx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	return eg.Wait()
}
