package messagesrepo

import (
	"time"

	"github.com/FischukSergey/chat-service/internal/store"
	"github.com/FischukSergey/chat-service/internal/types"
)

type Message struct {
	ID       types.MessageID
	ChatID   types.ChatID
	AuthorID types.UserID
	Body     string
	// FIXME: Остальные поля (тесты подскажут)
	IsVisibleForClient bool
	IsVisibleForManager bool
	IsService bool
	IsBlocked bool
	CreatedAt time.Time
}

func adaptStoreMessage(m *store.Message) Message {
	return Message{
		ID:       m.ID,
		ChatID:   m.ChatID,
		AuthorID: m.AuthorID,
		Body:     m.Body,
		// FIXME: Остальные поля (тесты подскажут)
		IsVisibleForClient: m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		IsService: m.IsService,
		IsBlocked: m.IsBlocked,
		CreatedAt: m.CreatedAt,
	}
}
