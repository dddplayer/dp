package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/internal/domain/code"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"path"
)

type Arch struct {
	*valueobject.CodeHandler
	relationDigraph *RelationDigraph
	directory       *Directory
}

func (arc *Arch) ObjectHandler() code.Handler {
	return arc.CodeHandler
}

func (arc *Arch) BuildHexagon() error {
	if err := arc.buildDirectory(); err != nil {
		return err
	}

	if arc.directory.ArchDesignPattern() != arch.DesignPatternHexagon {
		return fmt.Errorf("%s structure is only supported now", arch.DesignPatternHexagon)
	}

	if err := arc.buildOriginGraph(); err != nil {
		return err
	}

	return nil
}

func (arc *Arch) StrategicGraph() (arch.Diagram, error) {
	if err := arc.BuildHexagon(); err != nil {
		return nil, err
	}

	g, err := arc.buildStrategicArchGraph()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (arc *Arch) buildStrategicArchGraph() (*Diagram, error) {
	dm, err := NewDomainModel(arc.ObjRepo, arc.directory)
	if err != nil {
		return nil, err
	}
	if err := dm.StrategicGrouping(); err != nil {
		return nil, err
	}

	g, err := NewDiagram(arc.Scope, arch.PlainDiagram)
	if err != nil {
		return nil, err
	}

	for _, ag := range dm.aggregates {
		a, err := ag.Aggregate()
		if err != nil {
			return nil, err
		}
		if a.Entity == nil {
			fmt.Printf("warn: strategic aggregate %s has no entity", a.Name)
			continue
		}

		if err := g.AddObjTo(a, g.Name(), arch.RelationTypeAggregationRoot); err != nil {
			return nil, err
		}

		var sgObjs []arch.Object
		for _, sg := range ag.SubGroups() {
			switch sg.(type) {
			case *valueobject.EntityGroup:
				es := sg.(*valueobject.EntityGroup).Entities()
				for _, e := range es {
					if e.Identifier().ID() == a.Identifier().ID() {
						continue
					}
					sgObjs = append(sgObjs, e)
				}
			case *valueobject.VOGroup:
				vos := sg.(*valueobject.VOGroup).ValueObjects()
				sgObjs = append(sgObjs, vos.Objects()...)
			}
		}

		for _, o := range sgObjs {
			if err := g.AddObjTo(o, a.Identifier().ID(), arch.RelationTypeAggregation); err != nil {
				return nil, err
			}
		}
	}

	if err := arc.summaryDomainComponentRelations(g); err != nil {
		return nil, err
	}

	return g, nil
}

func (arc *Arch) summaryDomainComponentRelations(g *Diagram) error {
	combinations := generateCombinations(g.Objects())
	for _, comb := range combinations {
		if _, ok := comb.First.(arch.DomainObj); !ok {
			return fmt.Errorf("first object: %s is not a domain component", comb.First.Identifier().ID())
		}
		if _, ok := comb.Second.(arch.DomainObj); !ok {
			return fmt.Errorf("second object: %s is not a domain component", comb.Second.Identifier().ID())
		}

		metas, err := arc.relationDigraph.SummaryRelationMetas(
			comb.First.(arch.DomainObj).OriginIdentifier(),
			comb.Second.(arch.DomainObj).OriginIdentifier(),
		)
		if err != nil {
			return err
		}
		if err := g.AddRelations(comb.First.Identifier().ID(), comb.Second.Identifier().ID(), metas); err != nil {
			return err
		}
	}
	return nil
}

func (arc *Arch) TacticGraph(ops arch.Options) (arch.Diagram, error) {
	if err := arc.BuildHexagon(); err != nil {
		return nil, err
	}

	g, err := arc.buildTacticArchGraph()
	if err != nil {
		return nil, err
	}

	if ops.ShowAllRelations() {
		if err := arc.domainComponentRelations(g); err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (arc *Arch) buildTacticArchGraph() (*Diagram, error) {
	dm, err := NewDomainModel(arc.ObjRepo, arc.directory)
	if err != nil {
		return nil, err
	}
	if err := dm.TacticGrouping(); err != nil {
		return nil, err
	}

	g, err := NewDiagram(arc.Scope, arch.TableDiagram)
	if err != nil {
		return nil, err
	}

	for _, ag := range dm.aggregates {
		a, err := ag.Aggregate()
		if err != nil {
			return nil, err
		}
		if a.Entity == nil {
			fmt.Printf("warn: tactic aggregate %s has no entity", a.Name)
			continue
		}
		if err := dm.addAggregateToDiagram(g, ag, a); err != nil {
			return nil, err
		}

		for _, sg := range ag.SubGroups() {
			componentKey := path.Join(ag.Name(), sg.Name())
			if err := dm.addNodeToAggregate(g, componentKey, a); err != nil {
				return nil, err
			}

			switch sg.(type) {
			case *valueobject.EntityGroup:
				es := sg.(*valueobject.EntityGroup).Entities()
				for _, e := range es {
					if e.Identifier().ID() == a.Identifier().ID() {
						continue
					}
					if err := dm.addEntityToNode(g, componentKey, e); err != nil {
						return nil, err
					}
				}
			case *valueobject.VOGroup:
				vos := sg.(*valueobject.VOGroup).ValueObjects()
				for _, v := range vos {
					if err := dm.addVOToNode(g, componentKey, v); err != nil {
						return nil, err
					}
				}
			}

			if dg, ok := sg.(valueobject.DomainGroup); ok {
				if err := dm.buildDomainComponents(g, dg, componentKey); err != nil {
					return nil, err
				}
			}
		}
	}

	return g, nil
}

func (arc *Arch) domainComponentRelations(g *Diagram) error {
	combinations := generateCombinations(g.Objects())
	for _, comb := range combinations {
		if _, ok := comb.First.(arch.DomainObj); !ok {
			return fmt.Errorf("first object: %s is not a domain component", comb.First.Identifier().ID())
		}
		if _, ok := comb.Second.(arch.DomainObj); !ok {
			return fmt.Errorf("second object: %s is not a domain component", comb.Second.Identifier().ID())
		}

		metas, err := arc.relationDigraph.RelationMetas(
			comb.First.(arch.DomainObj).OriginIdentifier(),
			comb.Second.(arch.DomainObj).OriginIdentifier(),
		)
		if err != nil {
			return err
		}
		if err := g.AddRelations(comb.First.Identifier().ID(), comb.Second.Identifier().ID(), metas); err != nil {
			return err
		}
	}
	return nil
}

type Combination struct {
	First  arch.Object
	Second arch.Object
}

func generateCombinations(arr []arch.Object) []Combination {
	var combinations []Combination
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr); j++ {
			if i != j {
				combination := Combination{
					First:  arr[i],
					Second: arr[j],
				}
				combinations = append(combinations, combination)
			}
		}
	}
	return combinations
}

func (arc *Arch) buildDirectory() error {
	var objPaths []string
	dirMap := make(map[string][]arch.ObjIdentifier)

	allIds := arc.ObjRepo.All()
	for _, id := range allIds {
		objPaths = append(objPaths, id.ID())

		objDir := id.Dir()
		if _, ok := dirMap[objDir]; !ok {
			dirMap[objDir] = []arch.ObjIdentifier{}
		}
		dirMap[objDir] = append(dirMap[objDir], id)
	}
	arc.directory = NewDirectory(objPaths)
	for dir, objs := range dirMap {
		if err := arc.directory.AddObjs(dir, objs); err != nil {
			return err
		}
	}

	return nil
}

func (arc *Arch) buildOriginGraph() error {
	arc.relationDigraph = &RelationDigraph{
		Graph: directed.NewDirectedGraph(),
	}
	allIds := arc.ObjRepo.All()
	for _, id := range allIds {
		if err := arc.relationDigraph.AddObj(id); err != nil {
			return err
		}
	}

	arc.RelRepo.Walk(func(rel arch.Relation) error {
		return arc.relationDigraph.AddRelation(rel)
	})

	return nil
}

func (arc *Arch) GeneralGraph(ops arch.Options) (arch.Diagram, error) {
	if err := arc.BuildPlain(); err != nil {
		return nil, err
	}

	g, err := arc.buildGeneralArchGraph()
	if err != nil {
		return nil, err
	}

	if ops.ShowAllRelations() {
		if err := arc.componentRelations(g); err != nil {
			return nil, err
		}
	} else if ops.ShowStructEmbeddedRelations() {
		if err := arc.componentAssociationRelations(g); err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (arc *Arch) MessageFlowDiagram(startPath, endPath, modPath string) (arch.Diagram, error) {
	if err := arc.BuildPlain(); err != nil {
		return nil, err
	}

	mf := &MessageFlow{
		directory:       arc.directory,
		relationDigraph: arc.relationDigraph,
		objRepo:         arc.ObjRepo,
		mainPkgPath:     startPath,
		endPkgPath:      endPath,
		modulePath:      modPath,
	}

	return mf.buildDiagram()
}

func (arc *Arch) BuildPlain() error {
	if err := arc.buildDirectory(); err != nil {
		return err
	}

	if err := arc.buildOriginGraph(); err != nil {
		return err
	}

	return nil
}

func (arc *Arch) buildGeneralArchGraph() (*Diagram, error) {
	gm, err := NewGeneralModel(arc.ObjRepo, arc.directory)
	if err != nil {
		return nil, err
	}

	gm.Grouping()

	g, err := NewDiagram(arc.Scope, arch.TableDiagram)
	if err != nil {
		return nil, err
	}

	if err := gm.addRootGroupToDiagram(g); err != nil {
		return nil, err
	}

	return g, nil
}

func (arc *Arch) componentRelations(g *Diagram) error {
	combinations := generateCombinations(g.Objects())
	for _, comb := range combinations {
		metas, err := arc.relationDigraph.RelationMetas(
			comb.First.Identifier(),
			comb.Second.Identifier(),
		)
		if err != nil {
			return err
		}
		if err := g.AddRelations(comb.First.Identifier().ID(), comb.Second.Identifier().ID(), metas); err != nil {
			return err
		}
	}
	return nil
}

func (arc *Arch) componentAssociationRelations(g *Diagram) error {
	combinations := generateCombinations(g.Objects())
	for _, comb := range combinations {
		metas, err := arc.relationDigraph.RelationMetas(
			comb.First.Identifier(),
			comb.Second.Identifier(),
		)
		if err != nil {
			return err
		}
		if err := g.AddRelations(comb.First.Identifier().ID(), comb.Second.Identifier().ID(), arc.filterAssociationMetas(metas)); err != nil {
			return err
		}
	}
	return nil
}

func (arc *Arch) filterAssociationMetas(metas []arch.RelationMeta) []arch.RelationMeta {
	var filteredMetas []arch.RelationMeta
	for _, meta := range metas {
		switch meta.Type() {
		case arch.RelationTypeAssociationOneOne, arch.RelationTypeAssociationOneMany, arch.RelationTypeAssociation:
			filteredMetas = append(filteredMetas, meta)
		}
	}
	return filteredMetas
}
