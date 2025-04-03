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
	e.GET("/log/level", s.LogLevel)
	e.PUT("/log/level", s.SetLogLevel)
	index.addPage("/log/level", "Get log level")

	// обработка "/debug/pprof/" и связанных команд
	e.GET("/debug/pprof/", s.pprofIndex)
	e.GET("/debug/pprof/cmdline", s.pprofCmdline)
	e.GET("/debug/pprof/profile", s.pprofProfile)
	e.GET("/debug/pprof/symbol", s.pprofSymbol)
	e.GET("/debug/pprof/trace", s.pprofTrace)
	e.GET("/debug/pprof/goroutine", s.pprofGoroutine)
	e.GET("/debug/pprof/heap", s.pprofHeap)
	e.GET("/debug/pprof/threadcreate", s.pprofThreadcreate)
	e.GET("/debug/pprof/block", s.pprofBlock)
	e.GET("/debug/pprof/mutex", s.pprofMutex)

	index.addPage("/debug/pprof/", "Get pprof index")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-minute profile")

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

func (s *Server) LogLevel(eCtx echo.Context) error {
	// Вернуть текущий уровень логирования
	return eCtx.JSON(http.StatusOK, map[string]string{
		"level": logger.GetLevel(),
	})
}

func (s *Server) SetLogLevel(eCtx echo.Context) error {
	// Получить уровень логирования
	level := eCtx.FormValue("level")
	if level == "" {
		return eCtx.JSON(http.StatusBadRequest, "level is required")
	}
	// Установить новый глобальный уровень логирования
	if err := logger.SetLevel(level); err != nil {
		return eCtx.JSON(http.StatusBadRequest, err.Error())
	}
	// логируем изменение уровня логирования
	s.lg.Info("log level changed", zap.String("level", logger.GetLevel()))

	return eCtx.JSON(http.StatusOK, map[string]string{
		"level": logger.GetLevel(),
	})
}

func (s *Server) pprofIndex(eCtx echo.Context) error {
	pprof.Index(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofCmdline(eCtx echo.Context) error {
	pprof.Cmdline(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofProfile(eCtx echo.Context) error {
	pprof.Profile(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofSymbol(eCtx echo.Context) error {
	pprof.Symbol(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofTrace(eCtx echo.Context) error {
	pprof.Trace(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofGoroutine(eCtx echo.Context) error {
	pprof.Index(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofHeap(eCtx echo.Context) error {
	pprof.Index(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofThreadcreate(eCtx echo.Context) error {
	pprof.Index(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofBlock(eCtx echo.Context) error {
	pprof.Index(eCtx.Response().Writer, eCtx.Request())
	return nil
}

func (s *Server) pprofMutex(eCtx echo.Context) error {
	pprof.Index(eCtx.Response().Writer, eCtx.Request())
	return nil
}
