package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"golang.org/x/exp/slices"
)

type RelationDigraph struct {
	*directed.Graph
}

func (g *RelationDigraph) AddObj(id arch.ObjIdentifier) error {
	if err := g.AddNode(id.ID(), id); err != nil {
		return err
	}
	return nil
}

func (g *RelationDigraph) AddRelation(rel arch.Relation) error {
	switch rel.(type) {
	case arch.DependenceRelation:
		dep := rel.(arch.DependenceRelation)
		fromId := dep.From().Identifier().ID()
		toId := dep.DependsOn().Identifier().ID()
		if err := g.AddEdge(fromId, toId, rel.Type(), valueobject.NewRelationPos(
			dep.From().Position(), dep.DependsOn().Position())); err != nil {
			return err
		}
	case arch.CompositionRelation:
		comp := rel.(arch.CompositionRelation)
		fromId := comp.From().Identifier().ID()
		toId := comp.Child().Identifier().ID()
		if err := g.AddEdge(fromId, toId, rel.Type(), valueobject.NewRelationPos(
			comp.From().Position(), comp.Child().Position())); err != nil {
			return err
		}
	case arch.EmbeddingRelation:
		emb := rel.(arch.EmbeddingRelation)
		fromId := emb.From().Identifier().ID()
		toId := emb.Embedded().Identifier().ID()
		if err := g.AddEdge(fromId, toId, rel.Type(), valueobject.NewRelationPos(
			emb.From().Position(), emb.Embedded().Position())); err != nil {
			return err
		}
	case arch.ImplementationRelation:
		impl := rel.(arch.ImplementationRelation)
		fromId := impl.From().Identifier().ID()
		for _, ifc := range impl.Implements() {
			toId := ifc.Identifier().ID()
			if err := g.AddEdge(fromId, toId, rel.Type(), valueobject.NewRelationPos(
				impl.From().Position(), ifc.Position())); err != nil {
				return err
			}
		}
	case arch.AssociationRelation:
		assoc := rel.(arch.AssociationRelation)
		fromId := assoc.From().Identifier().ID()
		toId := assoc.Refer().Identifier().ID()
		if err := g.AddEdge(fromId, toId, assoc.AssociationType(), valueobject.NewRelationPos(
			assoc.From().Position(), assoc.Refer().Position())); err != nil {
			return err
		}
	}
	return nil
}

func (g *RelationDigraph) RelationMetas(from, to arch.ObjIdentifier) ([]arch.RelationMeta, error) {
	f := g.FindNodeByKey(from.ID())
	t := g.FindNodeByKey(to.ID())
	if f == nil || t == nil {
		return nil, fmt.Errorf("from or to node not found in Digraph")
	}

	fe := f.Edges
	feTo := make(map[string]arch.RelationMeta)
	for _, e := range fe {
		pos := e.Value.(arch.RelationPos)
		feTo[e.To.Value.(arch.ObjIdentifier).ID()] =
			valueobject.NewRelationMeta(e.Type.(arch.RelationType), pos.From(), pos.To())
	}

	tn := []*directed.Node{t}
	var toObjIds []string
	for _, n := range tn {
		toObjIds = append(toObjIds, n.Value.(arch.ObjIdentifier).ID())
	}

	var relations []arch.RelationMeta
	for k, v := range feTo {
		if slices.Contains(toObjIds, k) {
			relations = append(relations, v)
		}
	}

	return relations, nil
}

func (g *RelationDigraph) SummaryRelationMetas(from, to arch.ObjIdentifier) ([]arch.RelationMeta, error) {
	f := g.FindNodeByKey(from.ID())
	t := g.FindNodeByKey(to.ID())
	if f == nil || t == nil {
		return nil, fmt.Errorf("from or to node not found in Digraph")
	}

	fe, err := g.obtainEdgesFromNodeTree(f)
	if err != nil {
		return nil, err
	}

	feTo := make(map[string]arch.RelationMeta)
	for _, e := range fe {
		pos := e.Value.(arch.RelationPos)
		feTo[e.To.Value.(arch.ObjIdentifier).ID()] =
			valueobject.NewRelationMeta(e.Type.(arch.RelationType), pos.From(), pos.To())
	}

	tn, err := g.obtainNodesFromNodeTree(t)
	if err != nil {
		return nil, err
	}

	var toObjIds []string
	for _, n := range tn {
		toObjIds = append(toObjIds, n.Value.(arch.ObjIdentifier).ID())
	}

	var relations []arch.RelationMeta
	for k, v := range feTo {
		if slices.Contains(toObjIds, k) {
			relations = append(relations, v)
		}
	}

	return relations, nil
}

func (g *RelationDigraph) obtainEdgesFromNodeTree(node *directed.Node) ([]*directed.Edge, error) {
	var edges []*directed.Edge
	for _, e := range node.Edges {
		if t, ok := e.Type.(arch.RelationType); ok {
			switch t {
			case arch.RelationTypeEmbedding, arch.RelationTypeComposition:
				es, err := g.obtainEdgesFromNodeTree(e.To)
				if err != nil {
					return nil, err
				}
				edges = append(edges, es...)
			}
		}
		edges = append(edges, e)
	}

	return edges, nil
}

func (g *RelationDigraph) obtainNodesFromNodeTree(node *directed.Node) ([]*directed.Node, error) {
	nodes := []*directed.Node{node}
	for _, e := range node.Edges {
		if t, ok := e.Type.(arch.RelationType); ok {
			switch t {
			case arch.RelationTypeEmbedding, arch.RelationTypeComposition:
				ns, err := g.obtainNodesFromNodeTree(e.To)
				if err != nil {
					return nil, err
				}
				nodes = append(nodes, ns...)
			}
		}
	}

	return nodes, nil
}
