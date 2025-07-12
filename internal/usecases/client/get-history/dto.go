package gethistory

import (
	"errors"
	"time"

	"github.com/FischukSergey/chat-service/internal/types"
	"github.com/FischukSergey/chat-service/internal/validator"
)

type Request struct {
	ID       types.RequestID `validate:"required"`
	ClientID types.UserID    `validate:"required"`
	PageSize int             `validate:"omitempty,gte=10,lte=100"`
	Cursor   string          `validate:"omitempty,base64url"` // Фикс: Мотайте на ус приёмы валидации!
}

func (r Request) Validate() error {
	// Фикс: Нужно доработать метод согласно тестам, покрыв два хитрых случая:
	// Фикс:	- не указан ни курсор, ни размер страницы
	if r.PageSize == 0 && r.Cursor == "" {
		return errors.New("page_size and cursor are required")
	}
	if r.PageSize != 0 && r.Cursor != "" {
		return errors.New("page_size and cursor are required")
	}

	// Фикс:	- и наоборот: указан и курсор, и размер страницы
	//
	// Фикс: Сделать это можно как и дополнительным кодом в Validate(), так и через структурные теги
	// Фикс: https://github.com/go-playground/validator?tab=readme-ov-file#other
	return validator.Validator.Struct(r)
}

type Response struct {
	// Фикс: Заполнить (тесты помогут)
	Messages   []Message
	NextCursor string
}

type Message struct {
	// Фикс: Заполнить (тесты помогут)
	ID                  types.MessageID `json:"id"`
	AuthorID            types.UserID    `json:"authorId"`
	Body                string          `json:"body"`
	CreatedAt           time.Time       `json:"createdAt"`
	IsReceived          bool            `json:"isReceived"`
	IsBlocked           bool            `json:"isBlocked"`
	IsService           bool            `json:"isService"`
	IsVisibleForManager bool            `json:"isVisibleForManager"`
}
