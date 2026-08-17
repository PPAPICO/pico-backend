package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Favorite holds the schema definition for the Favorite entity.
type Favorite struct {
	ent.Schema
}

// Fields of the Favorite.
func (Favorite) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Favorite.
func (Favorite) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("favorites").
			Unique().
			Required(),

		edge.From("policy", GovernmentPolicy.Type).
			Ref("favorites").
			Unique().
			Required(),
	}
}
