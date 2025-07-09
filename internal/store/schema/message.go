package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/FischukSergey/chat-service/internal/types"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(types.MessageID{}).
			DefaultFunc(func() types.MessageID {
				return types.NewMessageID()
			}).
			Unique().
			Immutable(),
		field.String("body").
			NotEmpty(),
		field.String("author_id").
			GoType(types.UserID{}).
			Immutable().
			NotEmpty(),
		field.Bool("is_visible_for_client").
			Default(true),
		field.Bool("is_visible_for_manager").
			Default(true),
		field.Bool("is_blocked").
			Default(false),
		field.Bool("is_service").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.String("chat_id").
			GoType(types.ChatID{}).
			Immutable(),
		field.String("problem_id").
			GoType(types.ProblemID{}).
			Optional().
			Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).
			Ref("messages").
			Field("chat_id").
			Unique().
			Required().
			Immutable(),
		edge.From("problem", Problem.Type).
			Ref("messages").
			Field("problem_id").
			Unique().
			Immutable(),
	}
}

// Indexes of the Message.
func (Message) Indexes() []ent.Index {
	return []ent.Index{
		// 1. Основной индекс для пагинации по чату (сортировка по времени)
		index.Fields("chat_id", "created_at"),

		// 2. Индекс для поиска сообщений по проблеме
		index.Fields("problem_id", "created_at"),

		// 3. Индекс для поиска сообщений автора в чате
		index.Fields("chat_id", "author_id", "created_at"),

		// 4. Индекс для фильтрации видимых сообщений для клиента
		index.Fields("chat_id", "is_visible_for_client", "created_at"),

		// 5. Индекс для фильтрации видимых сообщений для менеджера
		index.Fields("chat_id", "is_visible_for_manager", "created_at"),

		// 6. Индекс для поиска заблокированных сообщений
		index.Fields("is_blocked", "created_at"),

		// 7. Индекс для поиска служебных сообщений
		index.Fields("chat_id", "is_service", "created_at"),

		// 8. Составной индекс для автора (для поиска всех сообщений пользователя)
		index.Fields("author_id", "created_at"),
	}
}
