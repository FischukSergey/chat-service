// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/FischukSergey/chat-service/internal/store/chat"
	"github.com/FischukSergey/chat-service/internal/store/message"
	"github.com/FischukSergey/chat-service/internal/store/problem"
	"github.com/FischukSergey/chat-service/internal/types"
)

// MessageCreate is the builder for creating a Message entity.
type MessageCreate struct {
	config
	mutation *MessageMutation
	hooks    []Hook
}

// SetBody sets the "body" field.
func (mc *MessageCreate) SetBody(s string) *MessageCreate {
	mc.mutation.SetBody(s)
	return mc
}

// SetAuthorID sets the "author_id" field.
func (mc *MessageCreate) SetAuthorID(ti types.UserID) *MessageCreate {
	mc.mutation.SetAuthorID(ti)
	return mc
}

// SetIsVisibleForClient sets the "is_visible_for_client" field.
func (mc *MessageCreate) SetIsVisibleForClient(b bool) *MessageCreate {
	mc.mutation.SetIsVisibleForClient(b)
	return mc
}

// SetNillableIsVisibleForClient sets the "is_visible_for_client" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsVisibleForClient(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsVisibleForClient(*b)
	}
	return mc
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (mc *MessageCreate) SetIsVisibleForManager(b bool) *MessageCreate {
	mc.mutation.SetIsVisibleForManager(b)
	return mc
}

// SetNillableIsVisibleForManager sets the "is_visible_for_manager" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsVisibleForManager(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsVisibleForManager(*b)
	}
	return mc
}

// SetIsBlocked sets the "is_blocked" field.
func (mc *MessageCreate) SetIsBlocked(b bool) *MessageCreate {
	mc.mutation.SetIsBlocked(b)
	return mc
}

// SetNillableIsBlocked sets the "is_blocked" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsBlocked(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsBlocked(*b)
	}
	return mc
}

// SetIsService sets the "is_service" field.
func (mc *MessageCreate) SetIsService(b bool) *MessageCreate {
	mc.mutation.SetIsService(b)
	return mc
}

// SetNillableIsService sets the "is_service" field if the given value is not nil.
func (mc *MessageCreate) SetNillableIsService(b *bool) *MessageCreate {
	if b != nil {
		mc.SetIsService(*b)
	}
	return mc
}

// SetCreatedAt sets the "created_at" field.
func (mc *MessageCreate) SetCreatedAt(t time.Time) *MessageCreate {
	mc.mutation.SetCreatedAt(t)
	return mc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (mc *MessageCreate) SetNillableCreatedAt(t *time.Time) *MessageCreate {
	if t != nil {
		mc.SetCreatedAt(*t)
	}
	return mc
}

// SetChatID sets the "chat_id" field.
func (mc *MessageCreate) SetChatID(ti types.ChatID) *MessageCreate {
	mc.mutation.SetChatID(ti)
	return mc
}

// SetProblemID sets the "problem_id" field.
func (mc *MessageCreate) SetProblemID(ti types.ProblemID) *MessageCreate {
	mc.mutation.SetProblemID(ti)
	return mc
}

// SetNillableProblemID sets the "problem_id" field if the given value is not nil.
func (mc *MessageCreate) SetNillableProblemID(ti *types.ProblemID) *MessageCreate {
	if ti != nil {
		mc.SetProblemID(*ti)
	}
	return mc
}

// SetID sets the "id" field.
func (mc *MessageCreate) SetID(ti types.MessageID) *MessageCreate {
	mc.mutation.SetID(ti)
	return mc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (mc *MessageCreate) SetNillableID(ti *types.MessageID) *MessageCreate {
	if ti != nil {
		mc.SetID(*ti)
	}
	return mc
}

// SetChat sets the "chat" edge to the Chat entity.
func (mc *MessageCreate) SetChat(c *Chat) *MessageCreate {
	return mc.SetChatID(c.ID)
}

// SetProblem sets the "problem" edge to the Problem entity.
func (mc *MessageCreate) SetProblem(p *Problem) *MessageCreate {
	return mc.SetProblemID(p.ID)
}

// Mutation returns the MessageMutation object of the builder.
func (mc *MessageCreate) Mutation() *MessageMutation {
	return mc.mutation
}

// Save creates the Message in the database.
func (mc *MessageCreate) Save(ctx context.Context) (*Message, error) {
	mc.defaults()
	return withHooks(ctx, mc.sqlSave, mc.mutation, mc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mc *MessageCreate) SaveX(ctx context.Context) *Message {
	v, err := mc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mc *MessageCreate) Exec(ctx context.Context) error {
	_, err := mc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mc *MessageCreate) ExecX(ctx context.Context) {
	if err := mc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (mc *MessageCreate) defaults() {
	if _, ok := mc.mutation.IsVisibleForClient(); !ok {
		v := message.DefaultIsVisibleForClient
		mc.mutation.SetIsVisibleForClient(v)
	}
	if _, ok := mc.mutation.IsVisibleForManager(); !ok {
		v := message.DefaultIsVisibleForManager
		mc.mutation.SetIsVisibleForManager(v)
	}
	if _, ok := mc.mutation.IsBlocked(); !ok {
		v := message.DefaultIsBlocked
		mc.mutation.SetIsBlocked(v)
	}
	if _, ok := mc.mutation.IsService(); !ok {
		v := message.DefaultIsService
		mc.mutation.SetIsService(v)
	}
	if _, ok := mc.mutation.CreatedAt(); !ok {
		v := message.DefaultCreatedAt()
		mc.mutation.SetCreatedAt(v)
	}
	if _, ok := mc.mutation.ID(); !ok {
		v := message.DefaultID()
		mc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mc *MessageCreate) check() error {
	if _, ok := mc.mutation.Body(); !ok {
		return &ValidationError{Name: "body", err: errors.New(`store: missing required field "Message.body"`)}
	}
	if v, ok := mc.mutation.Body(); ok {
		if err := message.BodyValidator(v); err != nil {
			return &ValidationError{Name: "body", err: fmt.Errorf(`store: validator failed for field "Message.body": %w`, err)}
		}
	}
	if _, ok := mc.mutation.AuthorID(); !ok {
		return &ValidationError{Name: "author_id", err: errors.New(`store: missing required field "Message.author_id"`)}
	}
	if v, ok := mc.mutation.AuthorID(); ok {
		if err := message.AuthorIDValidator(v.String()); err != nil {
			return &ValidationError{Name: "author_id", err: fmt.Errorf(`store: validator failed for field "Message.author_id": %w`, err)}
		}
	}
	if _, ok := mc.mutation.IsVisibleForClient(); !ok {
		return &ValidationError{Name: "is_visible_for_client", err: errors.New(`store: missing required field "Message.is_visible_for_client"`)}
	}
	if _, ok := mc.mutation.IsVisibleForManager(); !ok {
		return &ValidationError{Name: "is_visible_for_manager", err: errors.New(`store: missing required field "Message.is_visible_for_manager"`)}
	}
	if _, ok := mc.mutation.IsBlocked(); !ok {
		return &ValidationError{Name: "is_blocked", err: errors.New(`store: missing required field "Message.is_blocked"`)}
	}
	if _, ok := mc.mutation.IsService(); !ok {
		return &ValidationError{Name: "is_service", err: errors.New(`store: missing required field "Message.is_service"`)}
	}
	if _, ok := mc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`store: missing required field "Message.created_at"`)}
	}
	if _, ok := mc.mutation.ChatID(); !ok {
		return &ValidationError{Name: "chat_id", err: errors.New(`store: missing required field "Message.chat_id"`)}
	}
	if v, ok := mc.mutation.ChatID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "chat_id", err: fmt.Errorf(`store: validator failed for field "Message.chat_id": %w`, err)}
		}
	}
	if v, ok := mc.mutation.ProblemID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "problem_id", err: fmt.Errorf(`store: validator failed for field "Message.problem_id": %w`, err)}
		}
	}
	if v, ok := mc.mutation.ID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`store: validator failed for field "Message.id": %w`, err)}
		}
	}
	if len(mc.mutation.ChatIDs()) == 0 {
		return &ValidationError{Name: "chat", err: errors.New(`store: missing required edge "Message.chat"`)}
	}
	return nil
}

func (mc *MessageCreate) sqlSave(ctx context.Context) (*Message, error) {
	if err := mc.check(); err != nil {
		return nil, err
	}
	_node, _spec := mc.createSpec()
	if err := sqlgraph.CreateNode(ctx, mc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*types.MessageID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	mc.mutation.id = &_node.ID
	mc.mutation.done = true
	return _node, nil
}

func (mc *MessageCreate) createSpec() (*Message, *sqlgraph.CreateSpec) {
	var (
		_node = &Message{config: mc.config}
		_spec = sqlgraph.NewCreateSpec(message.Table, sqlgraph.NewFieldSpec(message.FieldID, field.TypeString))
	)
	if id, ok := mc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := mc.mutation.Body(); ok {
		_spec.SetField(message.FieldBody, field.TypeString, value)
		_node.Body = value
	}
	if value, ok := mc.mutation.AuthorID(); ok {
		_spec.SetField(message.FieldAuthorID, field.TypeString, value)
		_node.AuthorID = value
	}
	if value, ok := mc.mutation.IsVisibleForClient(); ok {
		_spec.SetField(message.FieldIsVisibleForClient, field.TypeBool, value)
		_node.IsVisibleForClient = value
	}
	if value, ok := mc.mutation.IsVisibleForManager(); ok {
		_spec.SetField(message.FieldIsVisibleForManager, field.TypeBool, value)
		_node.IsVisibleForManager = value
	}
	if value, ok := mc.mutation.IsBlocked(); ok {
		_spec.SetField(message.FieldIsBlocked, field.TypeBool, value)
		_node.IsBlocked = value
	}
	if value, ok := mc.mutation.IsService(); ok {
		_spec.SetField(message.FieldIsService, field.TypeBool, value)
		_node.IsService = value
	}
	if value, ok := mc.mutation.CreatedAt(); ok {
		_spec.SetField(message.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := mc.mutation.ChatIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ChatTable,
			Columns: []string{message.ChatColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chat.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ChatID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := mc.mutation.ProblemIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   message.ProblemTable,
			Columns: []string{message.ProblemColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ProblemID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MessageCreateBulk is the builder for creating many Message entities in bulk.
type MessageCreateBulk struct {
	config
	err      error
	builders []*MessageCreate
}

// Save creates the Message entities in the database.
func (mcb *MessageCreateBulk) Save(ctx context.Context) ([]*Message, error) {
	if mcb.err != nil {
		return nil, mcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(mcb.builders))
	nodes := make([]*Message, len(mcb.builders))
	mutators := make([]Mutator, len(mcb.builders))
	for i := range mcb.builders {
		func(i int, root context.Context) {
			builder := mcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MessageMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, mcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, mcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, mcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (mcb *MessageCreateBulk) SaveX(ctx context.Context) []*Message {
	v, err := mcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mcb *MessageCreateBulk) Exec(ctx context.Context) error {
	_, err := mcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mcb *MessageCreateBulk) ExecX(ctx context.Context) {
	if err := mcb.Exec(ctx); err != nil {
		panic(err)
	}
}
