package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/FischukSergey/chat-service/internal/buildinfo"
	"github.com/FischukSergey/chat-service/internal/logger"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr string `option:"mandatory" validate:"required,hostname_port"`
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

	// создание логгера c именем "server-debug" и уровнем как у глобального логгера
	lg := zap.L().Named("server-debug")

	// создание эхо-сервера
	e := echo.New()
	e.Use(middleware.Recover())

	// создание сервера в котором будет запущен эхо-сервер
	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
	index := newIndexPage()

	e.GET("/version", s.Version)
	index.addPage("/version", "Get build information")

	// обработка "/log/level"
	e.PUT("/log/level", echo.WrapHandler(logger.GlobalLevel))
	e.GET("/log/level", echo.WrapHandler(logger.GlobalLevel))
	index.addPage("/log/level", "Get log level")

	// добавляем ручку для тестирования ERROR логов
	e.GET("/debug/error", s.DebugError)
	index.addPage("/debug/error", "Debug Sentry error event")

	// обработка "/debug/pprof/" и связанных команд
	{
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		e.GET("/debug/pprof/*", echo.WrapHandler(pprofMux))
		index.addPage("/debug/pprof/", "Go std profiler")
		index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")
	}

	e.GET("/", index.handler)
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx) //nolint:contextcheck // graceful shutdown with new context
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}

func (s *Server) Version(eCtx echo.Context) error {
	// Вернуть информацию о сборке в формате JSON
	return eCtx.JSON(http.StatusOK, buildinfo.BuildInfo)
}

// DebugError - генерирует ERROR лог для тестирования.
func (s *Server) DebugError(eCtx echo.Context) error {
	// Получаем сообщение из query параметра или используем дефолтное
	message := eCtx.QueryParam("message")
	if message == "" {
		message = "Test error message"
	}

	// Логируем с уровнем ERROR
	s.lg.Error("DEBUG ERROR triggered",
		zap.String("message", message),
		//	zap.String("remote_ip", eCtx.RealIP()),
		//	zap.String("user_agent", eCtx.Request().UserAgent()),
	)

	return eCtx.JSON(http.StatusOK, map[string]string{
		"status":         "success",
		"message":        "ERROR log generated successfully",
		"logged_message": message,
	})
}
