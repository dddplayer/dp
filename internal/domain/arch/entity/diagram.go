package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"strings"
)

type Diagram struct {
	*directed.Graph
	root *directed.Node
	objs []arch.Object
	t    arch.DiagramType
}

func NewDiagram(name string, t arch.DiagramType) (*Diagram, error) {
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	g := &Diagram{
		Graph: directed.NewDirectedGraph(),
		objs:  []arch.Object{},
		t:     t,
	}
	if err := g.AddNode(name, valueobject.NewStringObj(name)); err != nil {
		return nil, err
	}
	g.root = g.FindNodeByKey(name)
	return g, nil
}

func (g *Diagram) Objects() []arch.Object {
	return g.objs
}

func (g *Diagram) AppendObject(objects ...arch.Object) {
	g.objs = append(g.objs, objects...)
}

func (g *Diagram) AddObj(obj arch.Object) error {
	if err := g.AddNode(obj.Identifier().ID(), obj); err != nil {
		return err
	}
	return nil
}

func (g *Diagram) AddObjTo(obj arch.Object, pid string, t arch.RelationType) error {
	if err := g.AddObj(obj); err != nil {
		return err
	}
	g.AppendObject(obj)

	if err := g.AddEdge(pid, obj.Identifier().ID(), t, nil); err != nil {
		return err
	}
	return nil
}

func (g *Diagram) AddStringTo(obj string, pid string, t arch.RelationType) error {
	strObj := valueobject.NewStringObj(obj)
	if err := g.AddObj(strObj); err != nil {
		return err
	}
	if err := g.AddEdge(pid, strObj.Identifier().ID(), t, nil); err != nil {
		return err
	}
	return nil
}

func (g *Diagram) AddRelations(fromId, toId string, metas []arch.RelationMeta) error {
	for _, rel := range metas {
		if err := g.AddEdge(fromId, toId, rel.Type(), rel.Position()); err != nil {
			return err
		}
	}
	return nil
}

func (g *Diagram) Name() string {
	return g.root.Value.(arch.Object).Identifier().Name()
}

func (g *Diagram) Type() arch.DiagramType {
	return g.t
}

func (g *Diagram) countDepth(count int, dn *directed.Node) int {
	for _, e := range dn.Edges {
		rt := e.Type.(arch.RelationType)
		if rt == arch.RelationTypeAggregationRoot {
			count++
			return g.countDepth(count, e.To)
		}
	}

	return count
}

func (g *Diagram) Edges() []arch.Edge {
	var edges []arch.Edge
	visitedNodes := make(map[*directed.Node]bool)

	merged := g.parseNodeEdge(g.root, visitedNodes).merge()
	for _, e := range merged {
		edges = append(edges, e)
	}

	return edges
}

func (g *Diagram) parseNodeEdge(dn *directed.Node, visited map[*directed.Node]bool) edges {
	var es edges

	visited[dn] = true
	for _, e := range dn.Edges {
		if !visited[e.To] {
			es = append(es, g.parseNodeEdge(e.To, visited)...)
		}

		if g.ignoreEdge(e) {
			continue
		}

		pos := valueobject.NewEmptyRelationPos()
		if val, ok := e.Value.(arch.RelationPos); ok {
			pos = val
		}
		es = append(es, newEdge(e.From.Key, e.To.Key, e.Type.(arch.RelationType), pos.From(), pos.To()))
	}

	return es
}

func (g *Diagram) ignoreEdge(edge *directed.Edge) bool {
	switch edge.Type.(arch.RelationType) {
	case arch.RelationTypeAttribution, arch.RelationTypeBehavior,
		arch.RelationTypeEmbedding, arch.RelationTypeComposition:
		return true
	}

	if _, ok := edge.From.Value.(*valueobject.StringObj); ok {
		if edge.Type.(arch.RelationType) == arch.RelationTypeAggregation {
			return true
		}
	}

	return false
}

func (g *Diagram) SubDiagrams() []arch.SubDiagram {
	var sds []arch.SubDiagram
	sds = append(sds, g.parseSubDiagrams(g.root))

	for _, e := range g.root.Edges {
		switch e.Type.(arch.RelationType) {
		case arch.RelationTypeAbstraction, arch.RelationTypeAggregationRoot:
			sds = append(sds, g.parseSubDiagrams(e.To))
		}
	}

	return sds
}

func (g *Diagram) parseSubDiagrams(dn *directed.Node) arch.SubDiagram {
	obj := dn.Value.(arch.Object)
	sd := &subDiagram{
		name:      obj.Identifier().ID(),
		nodes:     []arch.Node{},
		elements:  []arch.Element{},
		subGraphs: []arch.SubDiagram{},
	}

	switch obj.(type) {
	case *valueobject.Aggregate:
		if a, ok := obj.(*valueobject.Aggregate); ok {
			sd.name = a.Domain()
			n := &node{obj.Identifier().ID(), obj.Identifier().Name(), string(objColor(obj))}
			sd.nodes = append(sd.nodes, n)
			sd.elements = append(sd.elements, newElement(n, elementTypeClass))
		}
	case *valueobject.StringObj:
		n := &node{obj.Identifier().ID(), obj.Identifier().Name(), string(objColor(obj))}
		sd.nodes = append(sd.nodes, n)
		sd.elements = append(sd.elements, newElement(n, elementTypeGeneral))
	}

	for _, e := range dn.Edges {
		switch e.Type.(arch.RelationType) {
		case arch.RelationTypeAggregationRoot, arch.RelationTypeAbstraction:
			ssd := g.parseSubDiagrams(e.To)
			sd.subGraphs = append(sd.subGraphs, ssd)
		case arch.RelationTypeAggregation:
			toObj := e.To.Value.(arch.Object)
			n := &node{toObj.Identifier().ID(), toObj.Identifier().Name(), string(objColor(toObj))}
			sd.nodes = append(sd.nodes, n)
			switch toObj.(type) {
			case *valueobject.Entity, *valueobject.ValueObject, *valueobject.DomainInterface,
				*valueobject.Class, *valueobject.Interface:
				sd.elements = append(sd.elements, newElement(n, elementTypeClass))
				g.parseNode(e.To, sd)
			default:
				if e := sd.findElement(obj.Identifier().ID()); e != nil {
					ele := e.(*element)
					ele.addLeft(n)
				}
			}
		case arch.RelationTypeAttribution:
			g.parseClassEdge(e, sd, obj)
		case arch.RelationTypeBehavior:
			g.parseClassEdge(e, sd, obj)
		}
	}
	return sd
}

func (g *Diagram) parseNode(dn *directed.Node, sd *subDiagram) {
	p := dn.Value.(arch.Object)
	for _, e := range dn.Edges {
		g.parseClassEdge(e, sd, p)
	}
}

func (g *Diagram) parseClassEdge(e *directed.Edge, sd *subDiagram, p arch.Object) {
	toObj := e.To.Value.(arch.Object)
	name := toObj.Identifier().Name()
	if strings.HasPrefix(name, p.Identifier().Name()) && len(name) > len(p.Identifier().Name()) {
		name = name[len(p.Identifier().Name())+p.Identifier().NameSeparatorLength():]
	}

	switch e.Type.(arch.RelationType) {
	case arch.RelationTypeAttribution:
		n := &node{toObj.Identifier().ID(), name, string(objColor(toObj))}
		sd.nodes = append(sd.nodes, n)
		if e := sd.findElement(p.Identifier().ID()); e != nil {
			ele := e.(*element)
			ele.addRight(n)
		}
	case arch.RelationTypeBehavior:
		n := &node{toObj.Identifier().ID(), name, string(objColorWithParent(toObj, p))}
		sd.nodes = append(sd.nodes, n)
		if e := sd.findElement(p.Identifier().ID()); e != nil {
			ele := e.(*element)
			ele.addLeft(n)
		}
	}
}

type subDiagram struct {
	name      string
	nodes     []arch.Node
	elements  []arch.Element
	subGraphs []arch.SubDiagram
}

func (sd *subDiagram) Name() string                 { return sd.name }
func (sd *subDiagram) Nodes() []arch.Node           { return sd.nodes }
func (sd *subDiagram) SubGraphs() []arch.SubDiagram { return sd.subGraphs }
func (sd *subDiagram) Summary() []arch.Element      { return sd.elements }
func (sd *subDiagram) findElement(id string) arch.Element {
	for _, e := range sd.elements {
		if e.ID() == id {
			return e
		}
	}
	return nil
}

type node struct {
	id    string
	name  string
	color string
}

func (n *node) ID() string    { return n.id }
func (n *node) Name() string  { return n.name }
func (n *node) Color() string { return n.color }

type elementType string

const (
	elementTypeClass   elementType = "class"
	elementTypeGeneral elementType = "general"
)

type element struct {
	*node
	t        elementType
	children []arch.Nodes
}

func (e *element) Children() []arch.Nodes { return e.children }
func (e *element) addLeft(n *node) {
	switch e.t {
	case elementTypeClass, elementTypeGeneral:
		nodes := e.children[0]
		nodes = append(nodes, n)
		e.children[0] = nodes
	}
}
func (e *element) addRight(n *node) {
	switch e.t {
	case elementTypeClass:
		nodes := e.children[1]
		nodes = append(nodes, n)
		e.children[1] = nodes
	}
}

func newElement(n *node, t elementType) *element {
	e := &element{
		node:     n,
		t:        t,
		children: []arch.Nodes{},
	}
	switch t {
	case elementTypeClass:
		e.children = []arch.Nodes{{}, {}}
	case elementTypeGeneral:
		e.children = []arch.Nodes{{}}
	}
	return e
}

type edge struct {
	from         string
	to           string
	relationType arch.RelationType
	fromPos      arch.Position
	toPos        arch.Position
}

func newEdge(from, to string, relationType arch.RelationType, fromPos, toPos arch.Position) *edge {
	return &edge{
		from:         from,
		to:           to,
		relationType: relationType,
		fromPos:      fromPos,
		toPos:        toPos,
	}
}

func (e *edge) From() string {
	return e.from
}

func (e *edge) To() string {
	return e.to
}

func (e *edge) Type() arch.RelationType {
	return e.relationType
}

func (e *edge) Key() string {
	return fmt.Sprintf("%s-%s-%d", e.From(), e.To(), e.Type())
}

func (e *edge) pos() arch.RelationPos {
	if e.fromPos != nil && e.toPos != nil {
		return valueobject.NewRelationPos(e.fromPos, e.toPos)
	}
	return nil
}

type mergedEdge struct {
	*edge
	count int
	pos   []arch.RelationPos
}

func (e *mergedEdge) Pos() []arch.RelationPos {
	return e.pos
}

func (e *mergedEdge) Count() int {
	return e.count
}

type edges []*edge

func (es edges) merge() []*mergedEdge {
	mergedEdges := make(map[string]*mergedEdge)
	for _, e := range es {
		key := e.Key()

		relPos := e.pos()
		if me, ok := mergedEdges[key]; ok {
			me.count++
			if relPos != nil {
				for _, existingPos := range me.pos {
					if !existingPos.From().IsEqual(relPos.From()) || !existingPos.To().IsEqual(relPos.To()) {
						me.pos = append(me.pos, relPos)
						break
					}
				}
			}
		} else {
			m := &mergedEdge{
				edge:  e,
				count: 1,
				pos:   []arch.RelationPos{},
			}
			if relPos != nil {
				m.pos = append(m.pos, relPos)
			}
			mergedEdges[key] = m
		}
	}

	mergedEdgesArray := make([]*mergedEdge, 0, len(mergedEdges))
	for _, e := range mergedEdges {
		mergedEdgesArray = append(mergedEdgesArray, e)
	}

	return mergedEdgesArray
}
