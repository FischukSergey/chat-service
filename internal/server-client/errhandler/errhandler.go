package errhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalerrors "github.com/FischukSergey/chat-service/internal/errors"
)

var _ echo.HTTPErrorHandler = Handler{}.Handle

//go:generate options-gen -out-filename=errhandler_options.gen.go -from-struct=Options
type Options struct {
	logger          *zap.Logger                                    `option:"mandatory" validate:"required"`
	productionMode  bool                                           `option:"mandatory"`
	responseBuilder func(code int, msg string, details string) any `option:"mandatory" validate:"required"`
}

type Handler struct {
	lg              *zap.Logger
	productionMode  bool
	responseBuilder func(code int, msg string, details string) any
}

func New(opts Options) (Handler, error) {
	// Фикс: Validation & construction
	if err := opts.Validate(); err != nil {
		return Handler{}, err
	}
	return Handler{
		lg:              opts.logger,
		productionMode:  opts.productionMode,
		responseBuilder: opts.responseBuilder,
	}, nil
}

func (h Handler) Handle(err error, eCtx echo.Context) {
	// Фикс: 1) Обработать ошибку с помощью internal/errors.ProcessServerError
	code, msg, details := internalerrors.ProcessServerError(err)
	// Фикс: 2) Вырезать детали, если включен h.productionMode
	if h.productionMode {
		details = ""
	}
	// Фикс: 3) Собрать ответ с помощью h.responseBuilder. Статус ответа всегда 200, настоящий код в теле ответа.
	response := h.responseBuilder(code, msg, details)
	// Фикс: 4) Вернуть ответ в формате JSON
	err = eCtx.JSON(http.StatusOK, response)
	if err != nil {
		h.lg.Error("error", zap.Error(err))
	}
}
