package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// NewRequestLogger создает middleware для логирования HTTP запросов.
// Игнорирует OPTIONS запросы.
// Логирует: latency, remote_ip, host, method, path, request_id, user_agent и status.
// Если в контексте есть user_id, то логирует и его.
// Если последующий в цепочке хендлер вернул ошибку, то добавляет в лог и её.
func NewRequestLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// Игнорируем OPTIONS запросы
			if req.Method == http.MethodOptions {
				return next(c)
			}

			start := time.Now()
			var err error

			// Запускаем следующий обработчик
			if err = next(c); err != nil {
				c.Error(err)
			}

			// Рассчитываем время выполнения запроса
			latency := time.Since(start)

			// Получаем request_id из контекста
			requestID := ""
			if reqIDVal := c.Get("request_id"); reqIDVal != nil {
				if reqIDStr, ok := reqIDVal.(string); ok {
					requestID = reqIDStr
				}
			}

			// Формируем базовые поля для лога
			fields := []zap.Field{
				zap.Duration("latency", latency),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("request_id", requestID),
				zap.String("user_agent", req.UserAgent()),
				zap.Int("status", res.Status),
			}

			// Добавляем user_id из контекста, если есть
			userID, ok := userID(c)
			if ok {
				fields = append(fields, zap.String("user_id", userID.String()))
			} else {
				fields = append(fields, zap.String("user_id", ""))
			}

			// Если была ошибка, добавляем её в лог
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("request completed with error", fields...)
			} else {
				logger.Info("request completed", fields...)
			}

			return err
		}
	}
}
