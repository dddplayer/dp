package valueobject

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/code"
	"strings"
)

type CodeHandler struct {
	Scope   string
	ObjRepo repository.ObjectRepository
	RelRepo repository.RelationRepository
	errors  []error
}

func (ch *CodeHandler) pushError(err error) {
	ch.errors = append(ch.errors, err)
}

func (ch *CodeHandler) NodeHandler(node *code.Node) {
	id := newIdentifier(node.Meta)
	pos := newPosition(node.Pos)

	switch node.Type {
	case code.TypeGenIdent, code.TypeGenFunc, code.TypeGenArray, code.TypeGenMap:
		ch.handleGenObj(id, pos)
	case code.TypeGenStruct:
		ch.handleClass(id, pos)
	case code.TypeGenStructField, code.TypeGenStructEmbeddedField:
		if node.Parent == nil {
			ch.pushError(fmt.Errorf("struct field:%s without parent", id.name))
			return
		}
		ch.handleAttribute(id, pos, newIdentifier(node.Parent.Meta))
	case code.TypeGenInterface:
		ch.handleInterface(id, pos)
	case code.TypeGenInterfaceMethod:
		if node.Parent == nil {
			ch.pushError(fmt.Errorf("interface method:%s without parent", id.name))
			return
		}
		ch.handleInterfaceMethod(id, pos, newIdentifier(node.Parent.Meta))
	case code.TypeFunc:
		if node.Parent != nil {
			ch.handleFunc(id, pos, newIdentifier(node.Parent.Meta), newPosition(node.Parent.Pos))
		} else {
			ch.handleFunc(id, pos, nil, nil)
		}
	default:
		ch.handleGenObj(id, pos)
	}
}

func (ch *CodeHandler) handleClass(id *ident, pos *pos) {
	c := &Class{
		obj:     &obj{id: id, pos: pos},
		attrs:   []*ident{},
		methods: []*ident{},
	}

	if objR := ch.ObjRepo.Find(c.Identifier()); objR != nil {
		if o, ok := objR.(*MissingReceiver); ok {
			for _, m := range o.methods {
				c.AppendMethod(m)
			}
		}
	}

	if err := ch.ObjRepo.Insert(c); err != nil {
		ch.pushError(err)
	}
}

func (ch *CodeHandler) handleAttribute(id *ident, pos *pos, pid *ident) {
	attr := &Attr{&obj{id: id, pos: pos}}
	if err := ch.ObjRepo.Insert(attr); err != nil {
		ch.pushError(err)
		return
	}
	if obj := ch.ObjRepo.Find(pid); obj != nil {
		if o, ok := obj.(*Class); ok {
			o.AppendAttribute(attr.id)
		}
	}
}

func (ch *CodeHandler) handleGenObj(id *ident, pos *pos) {
	genObj := &General{&obj{id: id, pos: pos}}
	if err := ch.ObjRepo.Insert(genObj); err != nil {
		ch.pushError(err)
	}
}

func (ch *CodeHandler) handleFunc(id *ident, pos *pos, pid *ident, parentPos *pos) {
	genObj := &Function{
		obj:      &obj{id: id, pos: pos},
		Receiver: pid,
	}
	if err := ch.ObjRepo.Insert(genObj); err != nil {
		ch.pushError(err)
		return
	}

	if pid != nil {
		if objR := ch.ObjRepo.Find(pid); objR != nil {
			if o, ok := objR.(*Class); ok {
				o.AppendMethod(id)
			} else if o, ok := objR.(*MissingReceiver); ok {
				o.AppendMethod(id)
			}

		} else {
			missingRec := &MissingReceiver{
				obj: &obj{id: pid, pos: parentPos},
				methods: []*ident{
					id,
				},
			}
			if err := ch.ObjRepo.Insert(missingRec); err != nil {
				ch.pushError(err)
			}
		}
	}
}

func (ch *CodeHandler) handleInterface(id *ident, pos *pos) {
	i := &Interface{
		obj: &obj{
			id:  id,
			pos: pos,
		},
		methods: []*ident{},
	}
	if err := ch.ObjRepo.Insert(i); err != nil {
		ch.pushError(err)
	}
}

func (ch *CodeHandler) handleInterfaceMethod(id *ident, pos *pos, pid *ident) {
	im := &InterfaceMethod{
		obj: &obj{
			id:  id,
			pos: pos,
		},
	}
	if err := ch.ObjRepo.Insert(im); err != nil {
		ch.pushError(err)
		return
	}
	if obj := ch.ObjRepo.Find(pid); obj != nil {
		if i, ok := obj.(*Interface); ok {
			i.Append(im)
		}
	}
}

func (ch *CodeHandler) LinkHandler(link *code.Link) {
	if strings.Contains(link.From.Meta.Pkg(), ch.Scope) == false ||
		strings.Contains(link.To.Meta.Pkg(), ch.Scope) == false {
		return
	}

	fromId := newIdentifier(link.From.Meta)
	fromId.fixTmpName()
	fromPos := newPosition(link.From.Pos)

	toId := newIdentifier(link.To.Meta)
	toId.fixTmpName()
	toPos := emptyPosition()
	if link.To.Pos != nil {
		toPos = newPosition(link.To.Pos)
	}

	var r arch.Relation

	switch link.From.Type | link.To.Type {
	case code.TypeGenStructField | code.TypeAny,
		code.TypeGenStructEmbeddedField | code.TypeAny:
		r = NewAssociation(&obj{id: fromId, pos: fromPos}, &obj{id: toId, pos: toPos}, arch.RelationType(link.Relation))
	case code.TypeAny | code.TypeGenInterface:
		r = NewImplementation(&obj{id: fromId, pos: fromPos}, &obj{id: toId, pos: toPos})
	case code.TypeFunc | code.TypeFunc:
		r = NewDependence(&obj{id: fromId, pos: fromPos}, &obj{id: toId, pos: toPos})
	case code.TypeGenStruct | code.TypeGenStructField,
		code.TypeAny | code.TypeFunc,
		code.TypeGenInterface | code.TypeGenInterfaceMethod:
		r = NewComposition(&obj{id: fromId, pos: fromPos}, &obj{id: toId, pos: toPos})
	case code.TypeGenStruct | code.TypeGenStructEmbeddedField:
		r = NewEmbedding(&obj{id: fromId, pos: fromPos}, &obj{id: toId, pos: toPos})
	}

	if err := ch.RelRepo.Insert(r); err != nil {
		ch.pushError(err)
	}
}
