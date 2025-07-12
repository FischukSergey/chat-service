package gethistory

import (
	"context"
	"errors"

	cursorDecode "github.com/FischukSergey/chat-service/internal/cursor"
	messagesrepo "github.com/FischukSergey/chat-service/internal/repositories/messages"
	"github.com/FischukSergey/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=gethistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetClientChatMessages(
		ctx context.Context,
		clientID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	// Фикс: Validation & construction
	if err := opts.Validate(); err != nil {
		return UseCase{}, err
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	// Фикс: 1) Если запрос невалиден, то возвращаем ErrInvalidRequest
	if err := req.Validate(); err != nil {
		return Response{}, ErrInvalidRequest
	}

	// Фикс: 2) Если не вышло декодировать (decode) курсор, то возвращаем ErrInvalidCursor
	var cursor *messagesrepo.Cursor
	if req.Cursor != "" {
		var c messagesrepo.Cursor
		err := cursorDecode.Decode(req.Cursor, &c)
		if err != nil {
			return Response{}, ErrInvalidCursor
		}
		cursor = &c
	}

	// Фикс: 3) Запрашиваем сообщения из репозитория
	messages, nextCursor, err := u.Options.msgRepo.GetClientChatMessages(ctx, req.ClientID, req.PageSize, cursor)
	if err != nil {
		if errors.Is(err, messagesrepo.ErrInvalidCursor) {
			return Response{}, ErrInvalidCursor
		}
		return Response{}, err
	}

	// кодируем новый курсор
	var nextCursorEncode string
	if nextCursor != nil {
		nextCursorEncode, err = cursorDecode.Encode(nextCursor)
		if err != nil {
			return Response{}, err
		}
	}

	// Преобразуем []messagesrepo.Message в []Message
	resultMessages := make([]Message, 0, len(messages))
	for _, m := range messages {
		resultMessages = append(resultMessages, Message{
			ID:                  m.ID,
			AuthorID:            m.AuthorID,
			Body:                m.Body,
			CreatedAt:           m.CreatedAt,
			IsBlocked:           m.IsBlocked,
			IsVisibleForManager: m.IsVisibleForManager,
			IsService:           m.IsService,
			IsReceived:          m.IsVisibleForManager && !m.IsBlocked, // Логика IsReceived
		})
	}

	response := Response{
		NextCursor: nextCursorEncode,
		Messages:   resultMessages,
	}

	return response, nil
}
