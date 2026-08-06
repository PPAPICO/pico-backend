package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),

		field.String("name").
			NotEmpty(),

		field.String("email").
			Unique().
			NotEmpty(),

		field.String("password").
			NotEmpty(),

		field.Int("age"),

		field.String("region"),

		field.Enum("gender").
			Values("MALE", "FEMALE", "OTHER"),

		field.Bool("is_student").
			Default(false),

		field.Bool("is_youth").
			Default(false),

		field.Bool("is_pregnant").
			Optional().
			Nillable(),

		field.Bool("is_business").
			Optional().
			Nillable(),

		field.Strings("interests").
			Optional(),

		field.Bool("is_disabled").
			Default(false),

		field.Enum("nationality").
			Values("DOMESTIC", "FOREIGN"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("policy_matches", PolicyMatch.Type),
	}
}
