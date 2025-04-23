package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

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
