package schema

import (
	"github.com/google/uuid"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type PolicyMatch struct {
	ent.Schema
}

func (PolicyMatch) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.Enum("match").
			Values(
				"POSSIBLE",
				"UNCERTAIN",
				"IMPOSSIBLE",
			),
	}
}

func (PolicyMatch) Edges() []ent.Edge {
	return []ent.Edge{
		edge.
			From("policy", GovernmentPolicy.Type).
			Ref("matches").
			Unique().
			Required(),

		edge.
			From("user", User.Type).
			Ref("policy_matches").
			Unique().
			Required(),
	}
}
