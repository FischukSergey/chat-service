package clientv1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/FischukSergey/chat-service/internal/middlewares"
	// "github.com/FischukSergey/chat-service/internal/types".
	gethistory "github.com/FischukSergey/chat-service/internal/usecases/client/get-history"
	"github.com/FischukSergey/chat-service/pkg/pointer"
)

// var stub = MessagesPage{Messages: []Message{
// 	{
// 		AuthorId:  pointer.Ptr(types.NewUserID()),
// 		Body:      "Здравствуйте! Разберёмся.",
// 		CreatedAt: time.Now(),
// 		Id:        types.NewMessageID(),
// 	},
// 	{
// 		AuthorId: pointer.Ptr(
// 			types.MustParse[types.UserID]("bbc3fa26-2961-400b-beec-6fc56d509c36"),
// 		),
// 		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
// 		CreatedAt: time.Now().Add(-time.Minute),
// 		Id:        types.NewMessageID(),
// 	},
// }}

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	// Логируем входящий запрос
	zap.L().Info("received getHistory request",
		zap.String("requestID", params.XRequestID.String()))

	// Читаем параметры запроса (хотя в данном случае не используем)
	var req GetHistoryRequest
	if err := eCtx.Bind(&req); err != nil {
		zap.L().Error("failed to bind request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	// Фикс: 1) За-bind-ить входящий запрос
	// Фикс: 2) Вызвать соответствующий юзкейс
	request := gethistory.Request{
		ID:       params.XRequestID,
		ClientID: clientID,
		Cursor:   pointer.Indirect(req.Cursor),   // безопасно извлекаем значение или zero value
		PageSize: pointer.Indirect(req.PageSize), // аналогично
	}

	usecaseResp, err := h.getHistoryUseCase.Handle(ctx, request)
	if err != nil {
		// Фикс: 3) Обработать gethistory.ErrInvalidRequest и gethistory.ErrInvalidCursor
		if err == gethistory.ErrInvalidRequest || err == gethistory.ErrInvalidCursor {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Фикс: 4) Сформировать ответ, обрабатывая возможное отсутствие автора у сообщения
	messages := make([]Message, len(usecaseResp.Messages))
	for i, msg := range usecaseResp.Messages {
		messages[i] = Message{
			AuthorId:   pointer.PtrWithZeroAsNil(msg.AuthorID),
			Body:       msg.Body,
			CreatedAt:  msg.CreatedAt,
			Id:         msg.ID,
			IsBlocked:  msg.IsBlocked,
			IsReceived: msg.IsReceived,
			IsService:  msg.IsService,
		}
	}

	// Создаем кастомную структуру для backward compatibility
	customResponse := map[string]any{
		"data": map[string]any{
			"messages": messages,
			"next":     usecaseResp.NextCursor,
		},
	}

	return eCtx.JSON(http.StatusOK, customResponse)
}
