package entity

import (
	"fmt"
	ca "github.com/dddplayer/core/codeanalysis/entity"
	"github.com/dddplayer/core/valueobject"
	"strings"
)

type DomainModel struct {
	Name string
	Repo Repository
}

func (dm *DomainModel) NodeHandler(node *ca.Node) {
	id := valueobject.NewIdentifier(node.ID)
	pos := valueobject.NewPosition(node.Pos)

	switch node.Type {
	case ca.TypeGenIdent, ca.TypeGenFunc, ca.TypeGenArray:
		dm.handleGenObj(id, pos)
	case ca.TypeGenStruct:
		if IsEntity(id.Path) {
			e := NewEntity(id, pos)
			if err := dm.Repo.Insert(e); err != nil {
				fmt.Println(err)
			}
		} else if IsValueObject(id.Path) {
			vo := NewValueObject(id, pos)
			if err := dm.Repo.Insert(vo); err != nil {
				fmt.Println(err)
			}
		} else {
			dm.handleClass(id, pos)
		}
	case ca.TypeGenStructField:
		if node.Parent != nil {
			attr := &Attr{&obj{id: &id, pos: &pos}}
			if err := dm.Repo.Insert(attr); err != nil {
				fmt.Println(err)
			}
			pid := valueobject.NewIdentifier(node.Parent.ID)
			if obj := dm.Repo.Find(pid); obj != nil {
				if o, ok := obj.(DomainClass); ok {
					o.AppendAttribute(attr.id)
				}
			}
		}
	case ca.TypeGenInterface:
		i := NewInterface(id, pos)
		if err := dm.Repo.Insert(i); err != nil {
			fmt.Println(err)
		}
	case ca.TypeGenInterfaceMethod:
		if node.Parent != nil {
			im := NewInterfaceMethod(id, pos)
			if err := dm.Repo.Insert(im); err != nil {
				fmt.Println(err)
			}
			pid := valueobject.NewIdentifier(node.Parent.ID)
			if obj := dm.Repo.Find(pid); obj != nil {
				if i, ok := obj.(*Interface); ok {
					i.Append(im)
				}
			}
		}
	case ca.TypeFunc:
		if node.Parent != nil {
			dm.handleFunc(id, pos)
			pid := valueobject.NewIdentifier(node.Parent.ID)
			if obj := dm.Repo.Find(pid); obj != nil {
				if o, ok := obj.(DomainClass); ok {
					o.AppendCommand(&id)
				}
			}
		} else {
			if IsFactory(id.Path) {
				f := NewFactory(id, pos)
				if err := dm.Repo.Insert(f); err != nil {
					fmt.Println(err)
				}
			} else if IsService(id.Path) {
				s := NewService(id, pos)
				if err := dm.Repo.Insert(s); err != nil {
					fmt.Println(err)
				}
			} else {
				dm.handleFunc(id, pos)
			}
		}
	default:
		dm.handleGenObj(id, pos)
	}
}

func (dm *DomainModel) handleClass(id valueobject.Identifier, pos valueobject.Position) {
	c := &Class{
		obj:     &obj{id: &id, pos: &pos},
		attrs:   []*valueobject.Identifier{},
		methods: []*valueobject.Identifier{},
	}

	if err := dm.Repo.Insert(c); err != nil {
		fmt.Println(err)
	}
}

func (dm *DomainModel) handleGenObj(id valueobject.Identifier, pos valueobject.Position) {
	genObj := &General{&obj{id: &id, pos: &pos}}
	if err := dm.Repo.Insert(genObj); err != nil {
		fmt.Println(err)
	}
}

func (dm *DomainModel) handleFunc(id valueobject.Identifier, pos valueobject.Position) {
	genObj := &Function{&obj{id: &id, pos: &pos}}
	if err := dm.Repo.Insert(genObj); err != nil {
		fmt.Println(err)
	}
}

func getName(name string) string {
	return strings.Split(name, "$")[0]
}

func (dm *DomainModel) LinkHandler(link *ca.Link) {
	if strings.Contains(link.From.ID.String(), dm.Name) == false ||
		strings.Contains(link.To.ID.String(), dm.Name) == false {
		return
	}

	fromId := &valueobject.Identifier{
		Path: link.From.ID.Path(),
		Name: getName(link.From.ID.Name()),
	}
	fromPos := valueobject.NewPosition(link.From.Pos)
	toId := &valueobject.Identifier{
		Path: link.To.ID.Path(),
		Name: getName(link.To.ID.Name()),
	}

	r := &valueobject.Relation{
		From: fromId,
		To:   toId,
		Pos:  &fromPos,
		Ship: valueobject.RelationShip(link.Relation),
	}
	switch link.From.Type {
	case ca.TypeGenStructField:
		r.Type = valueobject.TypeRefer
	case ca.TypeGenInterface:
		r.Type = valueobject.TypeImplements
	case ca.TypeFunc:
		r.Type = valueobject.TypeCall
	}

	if err := dm.Repo.Insert(r); err != nil {
		fmt.Println(err)
	}
}
