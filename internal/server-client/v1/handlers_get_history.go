package clientv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/FischukSergey/chat-service/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  types.NewUserID(),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  types.MustParse[types.UserID]("bbc3fa26-2961-400b-beec-6fc56d509c36"), // подставь ID своего пользователя
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	// Логируем входящий запрос
	zap.L().Info("received getHistory request",
		zap.String("requestID", params.XRequestID.String()))

	// Читаем параметры запроса (хотя в данном случае не используем)
	var req GetHistoryRequest
	if err := eCtx.Bind(&req); err != nil {
		zap.L().Error("failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	// Возвращаем stub-данные в формате JSON согласно контракту API
	response := GetHistoryResponse{
		Data: stub,
	}

	return eCtx.JSON(http.StatusOK, response)
}
