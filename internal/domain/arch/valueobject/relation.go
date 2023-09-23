package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
)

type relation struct {
	from    *obj
	to      *obj
	relType arch.RelationType
}

func (r *relation) Type() arch.RelationType { return r.relType }
func (r *relation) From() arch.Object       { return r.from }

type Dependence struct {
	*relation
}

func (d *Dependence) DependsOn() arch.Object {
	return d.to
}

func NewDependence(from, to *obj) arch.Relation {
	return &Dependence{
		relation: &relation{
			from:    from,
			to:      to,
			relType: arch.RelationTypeDependency,
		},
	}
}

type Composition struct {
	*relation
}

func (c *Composition) Child() arch.Object {
	return c.to
}

func NewComposition(from, to *obj) arch.Relation {
	return &Composition{
		relation: &relation{
			from:    from,
			to:      to,
			relType: arch.RelationTypeComposition,
		},
	}
}

type Embedding struct {
	*relation
}

func (e *Embedding) Embedded() arch.Object {
	return e.to
}

func NewEmbedding(from, to *obj) arch.Relation {
	return &Embedding{
		relation: &relation{
			from:    from,
			to:      to,
			relType: arch.RelationTypeEmbedding,
		},
	}
}

type Implementation struct {
	*relation
	to []arch.Object
}

func (i *Implementation) Implements() []arch.Object {
	return i.to
}

func (i *Implementation) Implemented(ifc arch.Object) {
	i.to = append(i.to, ifc)
}

func NewImplementation(from, to *obj) arch.Relation {
	impl := &Implementation{
		relation: &relation{
			from:    from,
			to:      nil,
			relType: arch.RelationTypeImplementation,
		},
		to: make([]arch.Object, 0),
	}
	impl.Implemented(to)
	return impl
}

type Association struct {
	*relation
	ship arch.RelationType
}

func (a *Association) Refer() arch.Object {
	return a.to
}
func (a *Association) AssociationType() arch.RelationType {
	return a.ship
}

func NewAssociation(from, to *obj, relationship arch.RelationType) arch.Relation {
	return &Association{
		relation: &relation{
			from:    from,
			to:      to,
			relType: arch.RelationTypeAssociation,
		},
		ship: relationship,
	}
}
