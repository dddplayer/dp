package factory

import (
	"fmt"
	valueobject2 "github.com/dddplayer/core/codeanalysis/valueobject"
	dot "github.com/dddplayer/core/dot/entity"
	"github.com/dddplayer/core/entity"
	"github.com/dddplayer/core/valueobject"
	"path"
	"strings"
)

type DotBuilder struct {
	Repo      entity.Repository
	Domain    *entity.Domain
	relations []*valueobject.Relation
}

func (db *DotBuilder) Build() (*valueobject.DotGraph, error) {
	db.Repo.Walk(func(obj entity.DomainObject) error {
		if r, ok := obj.(*valueobject.Relation); ok {
			db.relations = append(db.relations, r)
		} else {
			db.buildDomain(obj)
		}

		return nil
	})
	return db.buildGraph()
}

func (db *DotBuilder) buildDomain(obj entity.DomainObject) {
	d := db.Domain.GetDomain(getDomainPath(obj.Identifier().Path))
	if d != nil {
		if o, ok := obj.(entity.DomainComponent); ok {
			d.Components = append(d.Components, o)
		} else if i, ok := obj.(*entity.Interface); ok {
			d.Interfaces = append(d.Interfaces, i)
		} else {
			switch o := obj.(type) {
			case *entity.General:
				d.Components = append(d.Components, entity.NewGeneralComponent(o))
			case *entity.Function:
				d.Components = append(d.Components, entity.NewFunctionComponent(o))
			case *entity.Class:
				d.Components = append(d.Components, entity.NewClassComponent(o))
			}
		}
	}
}

func getDomainPath(objPath string) string {
	p := objPath
	for _, ck := range []entity.ComponentKind{entity.KindEntity, entity.KindValueObject,
		entity.KindFactory, entity.KindService, entity.KindRepository} {
		if strings.HasPrefix(path.Base(p), string(ck)) {
			return path.Dir(p)
		}
	}
	base := path.Base(p)
	if strings.Contains(base, DomainObjJoiner) {
		es := strings.Split(base, DomainObjJoiner)
		p = path.Join(path.Dir(p), es[0])
	}
	return p
}

func (db *DotBuilder) buildGraph() (*valueobject.DotGraph, error) {
	digraph := valueobject.NewDotGraph(db.Domain.Name)

	if err := db.buildNodes(db.Domain, digraph); err != nil {
		return nil, err
	}
	if err := db.buildEdges(digraph); err != nil {
		return nil, err
	}
	return digraph, nil
}

const (
	DomainPortJoiner = dot.DotPortJoiner
	DomainObjJoiner  = valueobject2.NodeJoiner
)

func port(id *valueobject.Identifier) string {
	domain := getDomainPath(id.Path)
	domainBase := path.Base(domain)

	if domain != id.Path {
		p := strings.TrimPrefix(id.Base(), domainBase)
		p = strings.TrimPrefix(p, DomainObjJoiner)
		return fmt.Sprintf("%s%s%s", domainBase, DomainPortJoiner, p)
	}
	return fmt.Sprintf("%s%s%s", domainBase, DomainPortJoiner, id.Base())
}

func (db *DotBuilder) buildEdges(g *valueobject.DotGraph) error {
	for _, r := range db.relations {
		g.AppendEdge(valueobject.NewDotEdge(port(r.From), port(r.To)))
	}
	return nil
}

func (db *DotBuilder) buildNodes(d *entity.Domain, g *valueobject.DotGraph) error {
	if n := db.buildNode(d); n != nil {
		g.AppendNode(n)
	}
	for _, sd := range d.SubDomains {
		if err := db.buildNodes(sd, g); err != nil {
			return err
		}
	}

	return nil
}

func (db *DotBuilder) buildNode(d *entity.Domain) dot.DotNode {
	n := valueobject.NewNode(path.Base(d.Name))

	services := valueobject.NewElement(fmt.Sprintf("%s_service", n.Name()), "service")
	factories := valueobject.NewElement(fmt.Sprintf("%s_factory", n.Name()), "factory")
	generals := valueobject.NewElement(fmt.Sprintf("%s_general", n.Name()), "general")
	functions := valueobject.NewElement(fmt.Sprintf("%s_function", n.Name()), "function")
	var keyObjs []dot.DotElement

	for _, com := range d.Components {
		e := valueobject.NewElement(com.Identifier().Base(), com.Identifier().Name)

		if ko, ok := com.(entity.DomainKeyObj); ok {
			var ms []dot.DotAttribute
			for _, mid := range ko.Commands() {
				ms = append(ms, valueobject.NewAttribute(mid.Base(), mid.Name, string(entity.ColorMethod)))
			}
			e.Append(ms)

			var as []dot.DotAttribute
			for _, aid := range ko.Attributes() {
				as = append(as, valueobject.NewAttribute(aid.Base(), aid.Name, string(entity.ColorAttribute)))
			}
			e.Append(as)

			switch com.Kind() {
			case entity.KindEntity:
				e.SetColor(string(entity.ColorEntity))
			case entity.KindValueObject:
				e.SetColor(string(entity.ColorValueObject))
			case entity.KindClass:
				e.SetColor(string(entity.ColorClass))
			}
			keyObjs = append(keyObjs, e)

		} else {
			switch o := com.(type) {
			case *entity.Service:
				services.Append(valueobject.NewAttribute(o.Identifier().Base(), o.Identifier().Name, string(entity.ColorService)))
			case *entity.Factory:
				factories.Append(valueobject.NewAttribute(o.Identifier().Base(), o.Identifier().Name, string(entity.ColorFactory)))
			case *entity.GeneralComponent:
				generals.Append(valueobject.NewAttribute(o.Identifier().Base(), o.Identifier().Name, string(entity.ColorGeneral)))
			case *entity.FunctionComponent:
				functions.Append(valueobject.NewAttribute(o.Identifier().Base(), o.Identifier().Name, string(entity.ColorFunc)))
			}
		}
	}

	if len(services.Attributes()) > 0 {
		n.Append(services)
	}
	if len(factories.Attributes()) > 0 {
		n.Append(factories)
	}
	n.Append(keyObjs...)
	if len(generals.Attributes()) > 0 {
		n.Append(generals)
	}
	if len(functions.Attributes()) > 0 {
		n.Append(functions)
	}

	for _, i := range d.Interfaces {
		e := valueobject.NewElement(i.Identifier().Base(), i.Identifier().Name)
		e.SetColor(string(entity.ColorInterface))

		var ms []dot.DotAttribute
		for _, m := range i.Methods {
			ms = append(ms, valueobject.NewAttribute(m.Base(), m.Name, string(entity.ColorMethod)))
		}
		e.Append(ms)
		n.Append(e)
	}

	return n
}
