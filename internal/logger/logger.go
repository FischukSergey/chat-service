package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/FischukSergey/chat-service/internal/buildinfo"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options -defaults-from=var
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	clock          zapcore.Clock
	dsnSentry      string `validate:"omitempty,url"`
	env            string `validate:"required,oneof=dev stage prod"`
}

// defaultOptions - стандартные опции для логгера.
// Используются, если пользователь не предоставил свои опции.
var defaultOptions = Options{
	clock: zapcore.DefaultClock, // Используем стандартные часы из zapcore
	env:   "dev",                // По умолчанию используем окружение dev
}

// GlobalLevel - глобальный уровень логирования.
var GlobalLevel zap.AtomicLevel

// SentryClient - клиент для отправки отчетов в Sentry.
var SentryClient *sentry.Client

// MustInit - инициализирует логгер с заданными опциями.
// Если опции не валидны, то функция вызовет panic.
func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

// Init - инициализирует логгер с заданными опциями.
// Если опции не валидны, то функция вернет ошибку.
func Init(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}

	// парсим log level.
	GlobalLevel = zap.NewAtomicLevel()

	switch opts.level {
	case "debug":
		GlobalLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		GlobalLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		GlobalLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		GlobalLevel.SetLevel(zapcore.ErrorLevel)
	}

	// настраиваем логер:
	// 	- ключ для имени Named-логера – "component"
	// 	- ключ для времени лога – "T", формат времени – ISO8601
	// 	- если включен productionMode,
	// 		то уровень кодируется заглавными буквами ("INFO"), формат вывода – JSON,
	// 		иначе уровень кодируется заглавными буквами, но с добавлением цвета; а формат вывода – console plain-text.

	// Реализация настройки логгера
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "component",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// выбираем формат вывода лога
	var encoder zapcore.Encoder
	if opts.productionMode {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// создаём новый zapcore.Core на базе STDOUT.
	cores := []zapcore.Core{
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			GlobalLevel,
		),
	}

	// Если указан DSN для Sentry, настраиваем интеграцию с Sentry
	if opts.dsnSentry != "" {
		// Инициализируем клиент Sentry
		var err error
		SentryClient, err = NewSentryClient(
			opts.dsnSentry,
			opts.env,
			buildinfo.BuildInfo.Main.Version,
		)
		if err != nil {
			return fmt.Errorf("failed to initialize Sentry client: %v", err)
		}
		// Настраиваем конфигурацию для Sentry
		cfg := zapsentry.Configuration{
			Level: zapcore.WarnLevel, // Отправляем в Sentry только логи уровня WARN и выше
			Tags: map[string]string{
				"component": "system",
			},
		}
		// Создаём новое ядро Sentry
		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(SentryClient))
		if err != nil {
			return fmt.Errorf("failed to initialize Sentry core: %v", err)
		}
		// Добавляем ядро Sentry к существующим ядрам
		cores = append(cores, core)
	}

	// создаём новый логгер
	l := zap.New(zapcore.NewTee(cores...), zap.WithClock(opts.clock))

	// заменяем глобальный логгер на новый
	zap.ReplaceGlobals(l)

	return nil
}

// Sync - синхронизирует логгер.
func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
