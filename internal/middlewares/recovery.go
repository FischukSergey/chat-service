package middlewares

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// NewRecovery создает middleware для восстановления после паники в запросах.
// Логирует случившуюся ошибку и стек вызовов.
func NewRecovery(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// Получаем стек вызовов
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					stack := buf[:n]

					// Форматируем стек вызовов
					formattedStack := strings.Builder{}
					stackStr := string(stack)
					lines := strings.Split(stackStr, "\n")
					for i, line := range lines {
						lineNum := i + 1
						if _, err := formattedStack.WriteString(fmt.Sprintf("%3d: %s\n", lineNum, line)); err != nil {
							logger.Error("failed to write stack trace", zap.Error(err))
						}
					}

					// Логируем ошибку и стек вызовов
					logger.Error("panic recovered",
						zap.Any("error", r),
						zap.String("stack_trace", formattedStack.String()),
						zap.String("url", c.Request().URL.String()),
						zap.String("method", c.Request().Method),
					)

					// Возвращаем 500 Internal Server Error
					err := c.JSON(500, map[string]any{
						"error": "Internal Server Error",
					})
					if err != nil {
						logger.Error("failed to send error response", zap.Error(err))
					}
				}
			}()

			return next(c)
		}
	}
}
