package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/FischukSergey/chat-service/internal/types"
)

// Chat holds the schema definition for the Chat entity.
type Chat struct {
	ent.Schema
}

// Fields of the Chat.
func (Chat) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(types.ChatID{}).
			DefaultFunc(func() types.ChatID {
				return types.NewChatID()
			}).
			Unique().
			Immutable(),
		field.String("client_id").
			GoType(types.UserID{}).
			Unique().
			Immutable().
			NotEmpty().
			Validate(func(id string) error {
				if id == types.UserIDNil.String() {
					return errors.New("client_id cannot be nil")
				}
				return nil
			}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Chat.
func (Chat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", Message.Type),
		edge.To("problems", Problem.Type),
	}
}
