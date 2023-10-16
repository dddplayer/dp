package valueobject

import "github.com/dddplayer/dp/internal/domain/arch"

type obj struct {
	id  *ident
	pos *pos
}

func (o *obj) Identifier() arch.ObjIdentifier { return o.id }
func (o *obj) Position() arch.Position        { return o.pos }

type General struct {
	*obj
}

type MissingReceiver struct {
	*obj
	methods []*ident
}

type Class struct {
	*obj
	attrs   []*ident
	methods []*ident
}

type Function struct {
	*obj
	Receiver *ident
}

type Attr struct {
	*obj
}

type Interface struct {
	*obj
	methods []*ident
}

type InterfaceMethod struct {
	*obj
}

func (k *Class) Attributes() []*ident { return k.attrs }
func (k *Class) AppendAttribute(id *ident) {
	k.attrs = append(k.attrs, id)
}
func (k *Class) Methods() []*ident { return k.methods }
func (k *Class) AppendMethod(id *ident) {
	k.methods = append(k.methods, id)
}

func (k *MissingReceiver) AppendMethod(id *ident) {
	k.methods = append(k.methods, id)
}

func (i *Interface) Append(m *InterfaceMethod) {
	i.methods = append(i.methods, m.id)
}
func (i *Interface) Methods() []*ident { return i.methods }

func NewClass(o arch.Object, as []arch.ObjIdentifier, ms []arch.ObjIdentifier) *Class {
	cla := &Class{
		obj:     NewObj(o),
		attrs:   make([]*ident, 0),
		methods: make([]*ident, 0),
	}

	for _, a := range as {
		cla.AppendAttribute(&ident{
			name: a.Name(),
			pkg:  a.Dir(),
		})
	}

	for _, m := range ms {
		cla.AppendMethod(&ident{
			name: m.Name(),
			pkg:  m.Dir(),
		})
	}
	return cla
}

func NewObj(o arch.Object) *obj {
	return &obj{
		id: &ident{
			name: o.Identifier().Name(),
			pkg:  o.Identifier().Dir(),
		},
		pos: &pos{
			filename: o.Position().Filename(),
			offset:   o.Position().Offset(),
			line:     o.Position().Line(),
			column:   o.Position().Column(),
		},
	}
}

func NewFunction(o arch.Object, r arch.ObjIdentifier) *Function {
	if r != nil {
		return &Function{
			obj: NewObj(o),
			Receiver: &ident{
				name: r.Name(),
				pkg:  r.Dir(),
			},
		}
	}
	return &Function{
		obj:      NewObj(o),
		Receiver: nil,
	}
}

func NewAttr(o arch.Object) *Attr {
	return &Attr{
		obj: NewObj(o),
	}
}

func NewInterface(o arch.Object, ms []arch.Object) *Interface {
	i := &Interface{
		obj:     NewObj(o),
		methods: make([]*ident, 0),
	}
	for _, m := range ms {
		i.Append(&InterfaceMethod{
			obj: NewObj(m),
		})
	}
	return i
}

func NewGeneral(o arch.Object) *General {
	return &General{
		obj: NewObj(o),
	}
}
