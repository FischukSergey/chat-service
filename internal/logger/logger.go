package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options -defaults-from=var
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
	clock          zapcore.Clock
}

// defaultOptions - стандартные опции для логгера.
var defaultOptions = Options{
	clock: zapcore.DefaultClock, // Используем стандартные часы из zapcore
}

// GlobalLevel - глобальный уровень логирования.
var GlobalLevel zap.AtomicLevel

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

	// создаём новый логгер
	l := zap.New(zapcore.NewTee(cores...), zap.WithClock(opts.clock))

	// заменяем глобальный логгер на новый
	zap.ReplaceGlobals(l)

	return nil
}

/*
// SetLevel - устанавливает новый уровень логирования.
func SetLevel(level string) error {
	if err := GlobalLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("invalid level: %v", err)
	}
	zap.L().Core().Enabled(GlobalLevel.Level())
	return nil
}

// GetLevel - возвращает текущий уровень логирования.
func GetLevel() string {
	return GlobalLevel.Level().String()
}
*/

// Sync - синхронизирует логгер.
func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
