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
}

func (ch *CodeHandler) NodeHandler(node *code.Node) {
	id := newIdentifier(node.Meta)
	pos := newPosition(node.Pos)

	switch node.Type {
	case code.TypeGenIdent, code.TypeGenFunc, code.TypeGenArray:
		ch.handleGenObj(id, pos)
	case code.TypeGenStruct:
		ch.handleClass(id, pos)
	case code.TypeGenStructField, code.TypeGenStructEmbeddedField:
		if node.Parent == nil {
			fmt.Println("struct field no parent error")
		}
		ch.handleAttribute(id, pos, newIdentifier(node.Parent.Meta))
	case code.TypeGenInterface:
		ch.handleInterface(id, pos)
	case code.TypeGenInterfaceMethod:
		if node.Parent == nil {
			fmt.Println("interface method no parent error")
		}
		ch.handleInterfaceMethod(id, pos, newIdentifier(node.Parent.Meta))
	case code.TypeFunc:
		if node.Parent != nil {
			ch.handleFunc(id, pos, newIdentifier(node.Parent.Meta))
		} else {
			ch.handleFunc(id, pos, nil)
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

	if err := ch.ObjRepo.Insert(c); err != nil {
		fmt.Println(err)
	}
}

func (ch *CodeHandler) handleAttribute(id *ident, pos *pos, pid *ident) {
	attr := &Attr{&obj{id: id, pos: pos}}
	if err := ch.ObjRepo.Insert(attr); err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}
}

func (ch *CodeHandler) handleFunc(id *ident, pos *pos, pid *ident) {
	genObj := &Function{
		obj:      &obj{id: id, pos: pos},
		Receiver: pid,
	}
	if err := ch.ObjRepo.Insert(genObj); err != nil {
		fmt.Println(err)
	}

	if pid != nil {
		if obj := ch.ObjRepo.Find(pid); obj != nil {
			if o, ok := obj.(*Class); ok {
				o.AppendMethod(id)
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
		fmt.Println(err)
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
		fmt.Println(err)
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
		fmt.Println(err)
	}
}