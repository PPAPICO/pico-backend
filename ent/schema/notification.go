package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.Int("type"),

		field.String("message"),

		field.JSON("metadata", map[string]any{}).
			Optional(),

		field.Bool("is_read").
			Default(false),

		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Notification.
func (Notification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("receiver", User.Type).
			Ref("notifications").
			Unique().
			Required(),
	}
}
