package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"path"
)

type GeneralModel struct {
	repo      repository.ObjectRepository
	directory *Directory
	rootGroup valueobject.Group
}

func NewGeneralModel(r repository.ObjectRepository, d *Directory) (*GeneralModel, error) {
	return &GeneralModel{
		repo:      r,
		directory: d,
		rootGroup: nil,
	}, nil
}

func (gm *GeneralModel) Grouping() error {
	rootDir := gm.directory.RootDir()

	gm.directory.WalkRootDir(func(dir string, objIds []arch.ObjIdentifier) error {
		objs, err := gm.repo.GetObjects(objIds)
		if err != nil {
			return err
		}

		if dir == rootDir {
			gm.rootGroup = valueobject.NewGroup(dir, objs...)
		} else {
			pg := gm.FindGroup(path.Dir(dir), gm.rootGroup)
			pg.AppendGroups(valueobject.NewGroup(dir, objs...))
		}

		return nil
	})

	return nil
}

func (gm *GeneralModel) FindGroup(name string, g valueobject.Group) valueobject.Group {
	if name == g.Name() {
		return g
	}
	for _, group := range g.SubGroups() {
		if gp := gm.FindGroup(name, group); gp != nil {
			return gp
		}
	}
	return nil
}

func (gm *GeneralModel) addRootGroupToDiagram(g *Diagram) error {
	if err := gm.buildComponents(g, gm.rootGroup); err != nil {
		return err
	}

	for _, subGroup := range gm.rootGroup.SubGroups() {
		if err := gm.addGroupToDiagram(g, subGroup, gm.rootGroup.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (gm *GeneralModel) addGroupToDiagram(g *Diagram, group valueobject.Group, pid string) error {
	if err := g.AddStringTo(group.Name(), pid, arch.RelationTypeAggregationRoot); err != nil {
		return err
	}
	if err := gm.buildComponents(g, group); err != nil {
		return err
	}

	for _, subGroup := range group.SubGroups() {
		if err := gm.addGroupToDiagram(g, subGroup, group.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (gm *GeneralModel) buildComponents(g *Diagram, group valueobject.Group) error {
	if err := gm.buildAbstractComponent(g, group, valueobject.GeneralComponent); err != nil {
		return err
	}

	if err := gm.buildAbstractComponent(g, group, valueobject.FunctionComponent); err != nil {
		return err
	}

	if err := gm.buildAttributeComponents(g, group, valueobject.ClassComponent); err != nil {
		return err
	}

	if err := gm.buildAttributeComponents(g, group, valueobject.InterfaceComponent); err != nil {
		return err
	}

	return nil
}

func (gm *GeneralModel) buildAbstractComponent(g *Diagram, group valueobject.Group, componentType valueobject.ComponentType) error {
	var objs []arch.Object
	componentKey := path.Join(group.Name(), string(componentType))

	switch componentType {
	case valueobject.GeneralComponent:
		generals := group.Generals()
		for _, general := range generals {
			objs = append(objs, general)
		}
	case valueobject.FunctionComponent:
		functions := group.Functions()
		for _, function := range functions {
			objs = append(objs, function)
		}
	default:
		return fmt.Errorf("unsupported objs type: %s", componentType)
	}

	if len(objs) > 0 {
		if err := g.AddStringTo(componentKey, group.Name(), arch.RelationTypeAbstraction); err != nil {
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

func (gm *GeneralModel) buildAttributeComponents(g *Diagram, group valueobject.Group, componentType valueobject.ComponentType) error {
	switch componentType {
	case valueobject.ClassComponent:
		classes := group.Classes()
		for _, cla := range classes {
			if err := g.AddObjTo(cla, group.Name(), arch.RelationTypeAggregation); err != nil {
				return err
			}
			if err := gm.addClass(g, cla); err != nil {
				return err
			}
		}
	case valueobject.InterfaceComponent:
		interfaces := group.Interfaces()
		for _, itf := range interfaces {
			if err := g.AddObjTo(itf, group.Name(), arch.RelationTypeAggregation); err != nil {
				return err
			}
			for _, m := range itf.Methods() {
				if mo := gm.repo.Find(m); mo != nil {
					if err := g.AddObjTo(mo, itf.Identifier().ID(), arch.RelationTypeBehavior); err != nil {
						return err
					}
				}

			}
		}
	default:
		return fmt.Errorf("unsupported objs type: %s", componentType)
	}

	return nil
}

func (gm *GeneralModel) addClass(g *Diagram, cla *valueobject.Class) error {
	for _, m := range cla.Methods() {
		if mo := gm.repo.Find(m); mo != nil {
			if err := g.AddObjTo(mo, cla.Identifier().ID(), arch.RelationTypeBehavior); err != nil {
				return err
			}
		}
	}
	for _, a := range cla.Attributes() {
		if ao := gm.repo.Find(a); ao != nil {
			if err := g.AddObjTo(ao, cla.Identifier().ID(), arch.RelationTypeAttribution); err != nil {
				return err
			}
		}
	}
	return nil
}
