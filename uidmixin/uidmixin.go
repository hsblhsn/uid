package uidmixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/hsblhsn/uid"
)

func NewUID(prefix string) *UID {
	return &UID{prefix: prefix}
}

var _ ent.Mixin = (*UID)(nil)

// UID defines an ent UID that captures the uid prefix for a type.
type UID struct {
	mixin.Schema
	prefix string
}

// Fields provides the id field.
func (m UID) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(uid.ID("")).
			Unique().
			Immutable().
			DefaultFunc(func() uid.ID {
				id := uid.MustNew(m.prefix)
				return id
			}),
	}
}

// Annotations returns the annotations for a Mixin instance.
func (m UID) Annotations() []schema.Annotation {
	return []schema.Annotation{
		UIDAnnotation{Prefix: m.prefix, Length: len(m.prefix)},
	}
}

// UIDAnnotation captures the id prefix for a type.
type UIDAnnotation struct {
	Prefix string
	Length int
}

// Name implements the ent Annotation interface.
func (a UIDAnnotation) Name() string {
	return "UID"
}
