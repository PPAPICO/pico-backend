package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type GovernmentPolicy struct {
	ent.Schema
}

func (GovernmentPolicy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.String("title").
			NotEmpty(),

		field.Text("description"),

		field.Int("region_code"),

		field.Time("start_date"),

		field.Time("end_date"),

		field.String("address"),

		field.Float("latitude"),

		field.Float("longitude"),
	}
}

func (GovernmentPolicy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("matches", PolicyMatch.Type),
	}
}
