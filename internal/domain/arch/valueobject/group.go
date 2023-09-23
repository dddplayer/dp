package valueobject

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"strings"
)

type ComponentType string

const (
	GeneralComponent    ComponentType = "general"
	FunctionComponent   ComponentType = "function"
	ClassComponent      ComponentType = "class"
	InterfaceComponent  ComponentType = "interface"
	EntityComponent     ComponentType = "entity"
	VOComponent         ComponentType = "valueobject"
	RepositoryComponent ComponentType = "repository"
	FactoryComponent    ComponentType = "factory"
)

type Group interface {
	Name() string
	SubGroups() []Group
	AppendGroups(...Group)
	Objects() []arch.Object
	AppendObjects(...arch.Object)
	Classes() []*Class
	Generals() []*General
	Functions() []*Function
	Interfaces() []*Interface
}

type DomainGroup interface {
	Group
	Domain() string
	DomainGenerals() []*DomainGeneral
	DomainFunctions() []*DomainFunction
	DomainClasses() []*DomainClass
	DomainInterfaces() []*DomainInterface
}

type group struct {
	name      string
	subGroups []Group
	objs      []arch.Object
}

func NewGroup(name string, objs ...arch.Object) Group {
	return &group{
		name:      name,
		subGroups: []Group{},
		objs:      objs,
	}
}

func (g *group) Name() string {
	return g.name
}

func (g *group) SubGroups() []Group {
	return g.subGroups
}

func (g *group) AppendGroups(groups ...Group) {
	g.subGroups = append(g.subGroups, groups...)
}

func (g *group) Objects() []arch.Object {
	return g.objs
}

func (g *group) AppendObjects(objs ...arch.Object) {
	g.objs = append(g.objs, objs...)
}

func (g *group) Classes() []*Class {
	var cs []*Class
	for _, obj := range g.objs {
		if cla, ok := obj.(*Class); ok {
			cs = append(cs, cla)
		}
	}
	return cs
}

func (g *group) Generals() []*General {
	var gs []*General
	for _, obj := range g.objs {
		if g, ok := obj.(*General); ok {
			gs = append(gs, g)
		}
	}
	return gs
}

func (g *group) Functions() []*Function {
	var fs []*Function
	for _, obj := range g.objs {
		if f, ok := obj.(*Function); ok {
			if f.Receiver == nil {
				fs = append(fs, f)
			}
		}
	}
	return fs
}

func (g *group) Interfaces() []*Interface {
	var is []*Interface
	for _, obj := range g.objs {
		if i, ok := obj.(*Interface); ok {
			is = append(is, i)
		}
	}
	return is
}

type domainGroup struct {
	*group
	domain string
}

func (dg *domainGroup) Domain() string {
	return dg.domain
}

func (dg *domainGroup) DomainClasses() []*DomainClass {
	var dcs []*DomainClass
	for _, cla := range dg.Classes() {
		var as []*DomainAttr
		for _, attrId := range cla.Attributes() {
			if da := dg.getDomainAttr(attrId); da != nil {
				as = append(as, da)
			}
		}

		var fs []*DomainFunction
		for _, funcId := range cla.Methods() {
			if df := dg.getDomainFunction(funcId); df != nil {
				fs = append(fs, df)
			}
		}

		dcs = append(dcs, NewDomainClass(cla, dg.domain, as, fs))
	}
	return dcs
}

func (dg *domainGroup) DomainGenerals() []*DomainGeneral {
	var dgs []*DomainGeneral
	for _, gen := range dg.Generals() {
		dgs = append(dgs, &DomainGeneral{domainObj: &domainObj{obj: gen.obj, domain: dg.domain}})
	}
	return dgs
}

func (dg *domainGroup) DomainFunctions() []*DomainFunction {
	var dfs []*DomainFunction
	for _, f := range dg.Functions() {
		dfs = append(dfs, &DomainFunction{domainObj: &domainObj{obj: f.obj, domain: dg.domain}})
	}
	return dfs
}

func (dg *domainGroup) DomainInterfaces() []*DomainInterface {
	var difs []*DomainInterface
	for _, i := range dg.Interfaces() {

		var fs []*DomainFunction
		for _, funcId := range i.Methods() {
			if df := dg.getInterfaceFunction(funcId); df != nil {
				fs = append(fs, df)
			}
		}

		difs = append(difs, &DomainInterface{
			domainObj: &domainObj{obj: i.obj, domain: dg.domain},
			Methods:   fs,
		})
	}
	return difs
}

type EntityGroup struct {
	*domainGroup
}

func NewEntityGroup(domain string, classes ...arch.Object) *EntityGroup {
	return &EntityGroup{
		domainGroup: &domainGroup{
			domain: domain,
			group: &group{
				name:      string(EntityComponent),
				subGroups: []Group{},
				objs:      classes,
			},
		},
	}
}

func (dg *domainGroup) getDomainAttr(id *ident) *DomainAttr {
	var da *DomainAttr
	for _, o := range dg.objs {
		if o.Identifier().ID() == id.ID() {
			if cla, ok := o.(*Attr); ok {
				da = &DomainAttr{domainObj: &domainObj{obj: cla.obj, domain: dg.domain}}
			}
		}
	}
	return da
}

func (dg *domainGroup) getDomainFunction(id *ident) *DomainFunction {
	var df *DomainFunction
	for _, o := range dg.objs {
		if o.Identifier().ID() == id.ID() {
			if f, ok := o.(*Function); ok {
				if f.Receiver != nil {
					df = &DomainFunction{domainObj: &domainObj{obj: f.obj, domain: dg.domain}}
				}
			}
		}
	}
	return df
}

func (dg *domainGroup) getInterfaceFunction(id *ident) *DomainFunction {
	var df *DomainFunction
	for _, o := range dg.objs {
		if o.Identifier().ID() == id.ID() {
			if f, ok := o.(*InterfaceMethod); ok {
				df = &DomainFunction{domainObj: &domainObj{obj: f.obj, domain: dg.domain}}
			}
		}
	}
	return df
}

type Entities []*Entity

func (es Entities) Objects() []arch.Object {
	var objs []arch.Object
	for _, e := range es {
		objs = append(objs, e)
	}
	return objs
}

func (eg *EntityGroup) Entities() Entities {
	var es Entities
	for _, cla := range eg.DomainClasses() {
		es = append(es, &Entity{
			DomainClass: cla,
		})
	}
	return es
}

type VOGroup struct {
	*domainGroup
}

func NewVOGroup(domain string, vos ...arch.Object) *VOGroup {
	return &VOGroup{
		domainGroup: &domainGroup{
			domain: domain,
			group: &group{
				name:      string(VOComponent),
				subGroups: []Group{},
				objs:      vos,
			},
		},
	}
}

type ValueObjects []*ValueObject

func (vos ValueObjects) Objects() []arch.Object {
	var objs []arch.Object
	for _, vo := range vos {
		objs = append(objs, vo)
	}
	return objs
}

func (eg *VOGroup) ValueObjects() ValueObjects {
	var es ValueObjects
	for _, cla := range eg.DomainClasses() {
		es = append(es, &ValueObject{
			DomainClass: cla,
		})
	}
	return es
}

type AggregateGroup struct {
	*domainGroup
}

func NewAggregateGroup(a *Aggregate, domain string) *AggregateGroup {
	return &AggregateGroup{
		domainGroup: &domainGroup{
			domain: domain,
			group: &group{
				name:      a.Name,
				subGroups: []Group{},
				objs:      []arch.Object{a},
			},
		},
	}
}

func (agg *AggregateGroup) DomainName() string {
	return agg.Domain()
}

func (agg *AggregateGroup) Aggregate() (*Aggregate, error) {
	for _, o := range agg.objs {
		if a, ok := o.(*Aggregate); ok {
			if a.Entity == nil {
				for _, sg := range agg.subGroups {
					if g, ok := sg.(*EntityGroup); ok {
						entities := g.Entities()
						for _, e := range entities {
							if strings.ToLower(e.Identifier().Name()) == strings.ToLower(a.Name) {
								a.Entity = e
								break
							}
						}
					}
				}
			}

			return a, nil
		}
	}
	return nil, fmt.Errorf("aggregate %s not found", agg.Name())
}

func (agg *AggregateGroup) IsValid() bool {
	a, err := agg.Aggregate()
	if err != nil {
		return false
	}
	return a.Entity != nil
}
