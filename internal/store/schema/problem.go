package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/FischukSergey/chat-service/internal/types"
)

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(types.ProblemID{}).
			DefaultFunc(func() types.ProblemID {
				return types.NewProblemID()
			}).
			Unique().
			Immutable(),
		field.String("manager_id").
			GoType(types.UserID{}).
			NotEmpty().
			Immutable(),
		field.Enum("status").
			Values("open", "in_progress", "resolved", "closed").
			Default("open"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.String("chat_id").
			GoType(types.ChatID{}).
			Immutable(),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).
			Ref("problems").
			Field("chat_id").
			Unique().
			Required().
			Immutable(),
		edge.To("messages", Message.Type).
			Immutable(),
	}
}
