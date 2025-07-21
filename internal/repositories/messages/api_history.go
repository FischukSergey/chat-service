package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/FischukSergey/chat-service/internal/store"
	"github.com/FischukSergey/chat-service/internal/store/chat"
	"github.com/FischukSergey/chat-service/internal/store/message"
	"github.com/FischukSergey/chat-service/internal/types"
)

var (
	ErrInvalidPageSize         = errors.New("invalid page size")
	ErrInvalidCursor           = errors.New("invalid cursor")
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
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
	// Фикс: 1) Если указан pageSize, то валидируем, что он в пределах [10, 100] и используем его для запроса.
	if pageSize != 0 && (pageSize < 10 || pageSize > 100) {
		return nil, nil, fmt.Errorf("%w: pageSize must be between 10 and 100", ErrInvalidPageSize)
	}

	// Фикс: 2) Если указан cursor, то используем данные из него, предварительно валидируя:
	// Фикс: pageSize аналогично пункту выше и LastCreatedAt на заполненность в принципе.
	//
	if cursor != nil {
		if cursor.PageSize < 10 || cursor.PageSize > 100 {
			return nil, nil, fmt.Errorf("%w: pageSize must be between 10 and 100", ErrInvalidCursor)
		}
		if cursor.LastCreatedAt.IsZero() {
			return nil, nil, fmt.Errorf("%w: LastCreatedAt cannot be zero", ErrInvalidCursor)
		}
	}

	// Фикс: 3) Предполагается, что API на этом уровне используют верно и передают или pageSize или курсор.
	// Фикс: При желании можно обрабатывать обратную ситуацию и возвращать ошибку или даже паниковать.
	// Фикс: При этом не забыть добавить соответствующий тест!
	//
	hasCursor := cursor != nil
	hasPageSize := pageSize > 0

	// Должен быть передан либо pageSize, либо cursor
	if !hasCursor && !hasPageSize {
		return nil, nil, fmt.Errorf("%w: either pageSize or cursor must be provided", ErrInvalidPaginationParams)
	}

	// определяем размер запроса и время после которого нужно брать сообщения
	var queryPageSize int
	var queryAfter *time.Time

	if hasCursor {
		queryPageSize = cursor.PageSize
		queryAfter = &cursor.LastCreatedAt
	} else {
		queryPageSize = pageSize
	}

	// Фикс: 4) Возвращаем очередную страницу сообщений в соответствии с параметрами запроса.
	// Фикс: Первое сообщение является последним по времени своего создания (наиболее свежее).
	//
	// Получаем чат клиента
	clientChat, err := r.db.Chat(ctx).
		Query().
		Where(chat.ClientID(clientID)).
		Only(ctx)
	if err != nil {
		if !store.IsNotFound(err) {
			return nil, nil, fmt.Errorf("failed to get client chat: %w", err)
		}
		// Создаем чат если его нет
		clientChat, err = r.db.Chat(ctx).
			Create().
			SetClientID(clientID).
			Save(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create client chat: %w", err)
		}
	}
	// Строим запрос для получения сообщений
	query := r.db.Message(ctx).
		Query().
		Where(
			message.ChatID(clientChat.ID),    // 7) нужно доставать сообщения из клиентского чата
			message.IsVisibleForClient(true), // 7) нужно доставать сообщения, видимые клиенту
		).
		Order(message.ByCreatedAt(sql.OrderDesc())).
		Limit(queryPageSize + 1)

	// Если есть cursor, применяем фильтр по времени
	if queryAfter != nil {
		query = query.Where(message.CreatedAtLT(*queryAfter))
	}

	// Выполняем запрос
	storeMessages, err := query.All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Фикс: 5) Если впереди есть ещё страницы, то возвращаем курсор на следующую страницу, иначе nil.
	hasNextPage := len(storeMessages) > queryPageSize
	if hasNextPage {
		storeMessages = storeMessages[:len(storeMessages)-1]
	}
	// Преобразуем в API формат
	messages := make([]Message, len(storeMessages))
	for i, msg := range storeMessages {
		messages[i] = Message{
			ID:                  msg.ID,
			ChatID:              msg.ChatID,
			AuthorID:            msg.AuthorID,
			Body:                msg.Body,
			IsVisibleForClient:  msg.IsVisibleForClient,
			IsVisibleForManager: msg.IsVisibleForManager,
			IsService:           msg.IsService,
			IsBlocked:           msg.IsBlocked,
			CreatedAt:           msg.CreatedAt,
		}
	}

	// Формируем cursor для следующей страницы
	var nextCursor *Cursor
	if hasNextPage && len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		nextCursor = &Cursor{
			LastCreatedAt: lastMessage.CreatedAt,
			PageSize:      queryPageSize,
		}
	}
	// Фикс: 6) Пользуемся TEST_PSQL_DEBUG, чтобы понять, не превращает ли ent наш код в SQL-запрос, похожий на дичь.
	//
	// Фикс: 7) Отдельно обратите внимание на то, что
	// Фикс: - нужно доставать сообщения из клиентского чата (чужие чаты не должны попадать в выборку);
	// Фикс: - нужно доставать сообщения, видимые клиенту.
	return messages, nextCursor, nil
}
