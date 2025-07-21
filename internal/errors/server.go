package errors

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

// ServerError is used to return custom error codes to client.
type ServerError struct {
	Code    int
	Message string
	cause   error
}

func NewServerError(code int, msg string, err error) *ServerError {
	// Фикс: Создаем новый ServerError с переданными параметрами
	return &ServerError{
		Code:    code,
		Message: msg,
		cause:   err,
	}
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %v", s.Message, s.cause)
}

// Unwrap returns the underlying cause error.
func (s *ServerError) Unwrap() error {
	return s.cause
}

func GetServerErrorCode(err error) int {
	code, _, _ := ProcessServerError(err)
	return code
}

// ProcessServerError tries to retrieve from given error it's code, message and some details.
// For example, that fields can be used to build error response for client.
func ProcessServerError(err error) (code int, msg string, details string) {
	// Фикс: Если в цепочке ошибок есть ServerError или echo.HTTPError,
	// то достаём code и msg из них, детали равны результату работы метода Error().
	var serverError *ServerError
	var httpError *echo.HTTPError
	if errors.As(err, &serverError) {
		return serverError.Code, serverError.Message, serverError.Error()
	}
	if errors.As(err, &httpError) {
		return httpError.Code, httpError.Message.(string), httpError.Error()
	}

	// Фикс: 2) Если известных нам ошибок обнаружить не удалось, то возвращаем значения по умолчанию:
	// Фикс: 	 code = 500, message = "something went wrong", details = err.Error()
	return 500, "something went wrong", err.Error()
}
