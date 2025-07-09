package messagesrepo

import (
	"context"
	"errors"
	"time"

	"github.com/FischukSergey/chat-service/internal/types"
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	// FIXME: 1) Если указан pageSize, то валидируем, что он в пределах [10, 100] и используем его для запроса.
	if pageSize < 10 || pageSize > 100 {
		return nil, nil, ErrInvalidPageSize
	}

	// FIXME: 2) Если указан cursor, то используем данные из него, предварительно валидируя:
	// FIXME: pageSize аналогично пункту выше и LastCreatedAt на заполненность в принципе.
	//
	if cursor != nil {
		if cursor.PageSize < 10 || cursor.PageSize > 100 {
			return nil, nil, ErrInvalidPageSize
		}
		if cursor.LastCreatedAt.IsZero() {
			return nil, nil, ErrInvalidCursor
		}
	}


	// FIXME: 3) Предполагается, что API на этом уровне используют верно и передают или pageSize или курсор.
	// FIXME: При желании можно обрабатывать обратную ситуацию и возвращать ошибку или даже паниковать.
	// FIXME: При этом не забыть добавить соответствующий тест!
	//
	// FIXME: 4) Возвращаем очередную страницу сообщений в соответствии с параметрами запроса.
	// FIXME: Первое сообщение является последним по времени своего создания (наиболее свежее).
	//
	// FIXME: 5) Если впереди есть ещё страницы, то возвращаем курсор на следующую страницу, иначе nil.
	//
	// FIXME: 6) Пользуемся TEST_PSQL_DEBUG, чтобы понять, не превращает ли ent наш код в SQL-запрос, похожий на дичь.
	//
	// FIXME: 7) Отдельно обратите внимание на то, что
	// FIXME: - нужно доставать сообщения из клиентского чата (чужие чаты не должны попадать в выборку);
	// FIXME: - нужно доставать сообщения, видимые клиенту.
	return nil, nil, nil
}
