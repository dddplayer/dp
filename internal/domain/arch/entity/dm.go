package entity

import (
	"errors"
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"path"
)

type DomainModel struct {
	repo       repository.ObjectRepository
	directory  *Directory
	aggregates []*valueobject.AggregateGroup
}

func NewDomainModel(r repository.ObjectRepository, d *Directory) (*DomainModel, error) {
	return &DomainModel{
		repo:       r,
		directory:  d,
		aggregates: []*valueobject.AggregateGroup{},
	}, nil
}

func (dm *DomainModel) DomainName() (string, error) {
	domainDir, err := dm.directory.DomainDir()
	if err != nil {
		return "", err
	}
	return domainDir, nil
}

func (dm *DomainModel) TacticGrouping() error {
	domainDir, err := dm.directory.DomainDir()
	if err != nil {
		return err
	}

	dm.directory.WalkDir(domainDir, func(dir string, objIds []arch.ObjIdentifier) error {
		switch dm.directory.HexagonDirectory(dir) {
		case arch.HexagonDirectoryDomain:
			return nil
		case arch.HexagonDirectoryAggregate:
			name := path.Base(dir)
			ag := valueobject.NewAggregateGroup(&valueobject.Aggregate{
				Name: name,
			}, path.Join(domainDir, name))
			dm.aggregates = append(dm.aggregates, ag)
			objs, err := dm.repo.GetObjects(objIds)
			if err != nil {
				return err
			}
			ag.AppendObjects(objs...)
		case arch.HexagonDirectoryEntity:
			err := dm.processObjects(objIds, valueobject.EntityComponent, dir)
			if err != nil {
				return err
			}
		case arch.HexagonDirectoryValueObject:
			err := dm.processObjects(objIds, valueobject.VOComponent, dir)
			if err != nil {
				return err
			}
		case arch.HexagonDirectoryInvalid:
			err = errors.New("invalid hexagon domain directory")
			return err
		}
		return nil
	})

	return nil
}

func (dm *DomainModel) StrategicGrouping() error {
	domainDir, err := dm.directory.DomainDir()
	if err != nil {
		return err
	}

	dm.directory.WalkDir(domainDir, func(dir string, objIds []arch.ObjIdentifier) error {
		switch dm.directory.HexagonDirectory(dir) {
		case arch.HexagonDirectoryDomain:
			return nil
		case arch.HexagonDirectoryAggregate:
			name := path.Base(dir)
			dm.aggregates = append(dm.aggregates, valueobject.NewAggregateGroup(
				&valueobject.Aggregate{
					Name: name,
				},
				path.Join(domainDir, name)))
		case arch.HexagonDirectoryEntity:
			err := dm.processClasses(objIds, valueobject.EntityComponent, dir)
			if err != nil {
				return err
			}
		case arch.HexagonDirectoryValueObject:
			err := dm.processClasses(objIds, valueobject.VOComponent, dir)
			if err != nil {
				return err
			}
		case arch.HexagonDirectoryRepository:
			//todo
		case arch.HexagonDirectoryFactory:
			//todo
		case arch.HexagonDirectoryInvalid:
			err = errors.New("invalid hexagon domain directory")
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (dm *DomainModel) processObjects(objIds []arch.ObjIdentifier, groupType valueobject.ComponentType, dir string) error {
	ag := dm.FindAggregateGroup(path.Base(path.Dir(dir)))
	objs, err := dm.repo.GetObjects(objIds)
	if err != nil {
		return err
	}

	return dm.processComponent(ag, groupType, objs)
}

func (dm *DomainModel) processClasses(objIds []arch.ObjIdentifier, groupType valueobject.ComponentType, dir string) error {
	ag := dm.FindAggregateGroup(path.Base(path.Dir(dir)))
	classes, err := dm.getClass(objIds)
	if err != nil {
		return err
	}

	return dm.processComponent(ag, groupType, classes)
}

func (dm *DomainModel) processComponent(ag *valueobject.AggregateGroup, groupType valueobject.ComponentType, objects []arch.Object) error {
	switch groupType {
	case valueobject.EntityComponent:
		ag.AppendGroups(valueobject.NewEntityGroup(ag.Domain(), objects...))
	case valueobject.VOComponent:
		ag.AppendGroups(valueobject.NewVOGroup(ag.Domain(), objects...))
	}

	return nil
}

func (dm *DomainModel) getClass(objIds []arch.ObjIdentifier) ([]arch.Object, error) {
	objs, err := dm.repo.GetObjects(objIds)
	if err != nil {
		return nil, err
	}
	var classes []arch.Object
	for _, obj := range objs {
		if cla, ok := obj.(*valueobject.Class); ok {
			classes = append(classes, cla)
		}
	}
	return classes, nil
}

func (dm *DomainModel) FindAggregateGroup(name string) *valueobject.AggregateGroup {
	for _, ag := range dm.aggregates {
		if ag.Name() == name {
			return ag
		}
	}
	return nil
}

func (dm *DomainModel) addEntityToNode(g *Diagram, pid string, e *valueobject.Entity) error {
	if err := g.AddObjTo(e, pid, arch.RelationTypeAggregation); err != nil {
		return err
	}
	return dm.addDomainClass(g, e.DomainClass)
}

func (dm *DomainModel) addVOToNode(g *Diagram, pid string, vo *valueobject.ValueObject) error {
	if err := g.AddObjTo(vo, pid, arch.RelationTypeAggregation); err != nil {
		return err
	}
	return dm.addDomainClass(g, vo.DomainClass)
}

func (dm *DomainModel) addNodeToAggregate(g *Diagram, key string, a *valueobject.Aggregate) error {
	return g.AddStringTo(key, a.Identifier().ID(), arch.RelationTypeAggregationRoot)
}

func (dm *DomainModel) addAggregateToDiagram(g *Diagram, ag *valueobject.AggregateGroup, a *valueobject.Aggregate) error {
	if err := g.AddObjTo(a, g.Name(), arch.RelationTypeAggregationRoot); err != nil {
		return err
	}
	if err := dm.addDomainClass(g, a.DomainClass); err != nil {
		return err
	}

	if err := dm.buildDomainComponents(g, ag, a.Identifier().ID()); err != nil {
		return err
	}
	return nil
}

func (dm *DomainModel) addDomainClass(g *Diagram, cla *valueobject.DomainClass) error {
	for _, dm := range cla.Methods {
		if err := g.AddObjTo(dm, cla.Identifier().ID(), arch.RelationTypeBehavior); err != nil {
			return err
		}
	}
	for _, da := range cla.Attributes {
		if err := g.AddObjTo(da, cla.Identifier().ID(), arch.RelationTypeAttribution); err != nil {
			return err
		}
	}

	return nil
}

func (dm *DomainModel) buildDomainComponents(g *Diagram, group valueobject.DomainGroup, pid string) error {
	if err := dm.buildAbstractComponent(g, group, pid, valueobject.GeneralComponent); err != nil {
		return err
	}

	if err := dm.buildAbstractComponent(g, group, pid, valueobject.FunctionComponent); err != nil {
		return err
	}

	_, isEntityGroup := group.(*valueobject.EntityGroup)
	_, isVOGroup := group.(*valueobject.VOGroup)
	if !isEntityGroup && !isVOGroup {
		if err := dm.buildComponents(g, group, pid, valueobject.ClassComponent); err != nil {
			return err
		}
	}

	if err := dm.buildComponents(g, group, pid, valueobject.InterfaceComponent); err != nil {
		return err
	}

	return nil
}

func (dm *DomainModel) buildComponents(g *Diagram, group valueobject.DomainGroup, pid string, componentType valueobject.ComponentType) error {
	switch componentType {
	case valueobject.ClassComponent:
		classes := group.DomainClasses()
		for _, cla := range classes {
			if err := g.AddObjTo(cla, pid, arch.RelationTypeAggregation); err != nil {
				return err
			}
			if err := dm.addDomainClass(g, cla); err != nil {
				return err
			}
		}
	case valueobject.InterfaceComponent:
		interfaces := group.DomainInterfaces()
		for _, itf := range interfaces {
			if err := g.AddObjTo(itf, pid, arch.RelationTypeAggregation); err != nil {
				return err
			}
			for _, m := range itf.Methods {
				if err := g.AddObjTo(m, itf.Identifier().ID(), arch.RelationTypeBehavior); err != nil {
					return err
				}
			}
		}
	default:
		return fmt.Errorf("unsupported objs type: %s", componentType)
	}

	return nil
}

func (dm *DomainModel) buildAbstractComponent(g *Diagram, group valueobject.DomainGroup, pid string, componentType valueobject.ComponentType) error {
	var objs arch.DomainObjs
	componentKey := path.Join(pid, string(componentType))

	switch componentType {
	case valueobject.GeneralComponent:
		generals := group.DomainGenerals()
		for _, general := range generals {
			objs = append(objs, general)
		}
	case valueobject.FunctionComponent:
		functions := group.DomainFunctions()
		for _, function := range functions {
			objs = append(objs, function)
		}
	default:
		return fmt.Errorf("unsupported objs type: %s", componentType)
	}

	if len(objs) > 0 {
		if err := g.AddStringTo(componentKey, pid, arch.RelationTypeAbstraction); err != nil {
			return err
		}

		for _, o := range objs {
			if err := g.AddObjTo(o, componentKey, arch.RelationTypeAggregation); err != nil {
				return err
			}
		}
	}

	return nil
}
