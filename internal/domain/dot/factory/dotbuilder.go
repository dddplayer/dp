package factory

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/dot"
	"github.com/dddplayer/dp/internal/domain/dot/entity"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"path"
	"strconv"
	"strings"
)

func NewDotBuilder(diagram arch.Diagram) *DotBuilder {
	return &DotBuilder{
		archDiagram: diagram,
		portMap:     make(map[string]string),
	}
}

type DotBuilder struct {
	archDiagram arch.Diagram
	dot         *entity.Dot
	portMap     map[string]string
}

func (db *DotBuilder) Build() (*entity.Dot, error) {
	db.dot = &entity.Dot{
		Name:      db.archDiagram.Name(),
		Label:     fmt.Sprintf("\n\n%s\nDomain Model\n\nPowered by DDD Player", db.archDiagram.Name()),
		SubGraphs: []*entity.SubGraph{},
		Edges:     []*entity.Edge{},
	}

	if err := db.buildSubGraph(); err != nil {
		return nil, err
	}
	if err := db.buildEdges(); err != nil {
		return nil, err
	}
	if err := db.buildTemplates(); err != nil {
		return nil, err
	}

	return db.dot, nil
}

func (db *DotBuilder) buildSubGraph() error {
	var g *entity.SubGraph
	for _, sd := range db.archDiagram.SubDiagrams() {
		if db.isDeepMode() {
			g = NewSummarySubGraph(sd)
			db.buildPortMap(g)
		} else {
			g = NewSubGraph(sd)
		}
		db.dot.SubGraphs = append(db.dot.SubGraphs, g)

		for _, sg := range sd.SubGraphs() {
			if err := db.buildGraph(sg, g); err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *DotBuilder) buildGraph(sd arch.SubDiagram, pg *entity.SubGraph) error {
	var g *entity.SubGraph
	if db.isDeepMode() {
		g = NewSummarySubGraph(sd)
		db.buildPortMap(g)
	} else {
		g = NewSubGraph(sd)
	}
	pg.SubGraphs = append(pg.SubGraphs, g)

	for _, sg := range sd.SubGraphs() {
		if err := db.buildGraph(sg, g); err != nil {
			return err
		}
	}

	return nil
}

func (db *DotBuilder) buildPortMap(g *entity.SubGraph) {
	for _, n := range g.Nodes {
		for _, r := range n.Table.Rows {
			for _, d := range r.Data {
				if d.Port == "" {
					continue
				}
				db.portMap[d.Port] = n.ID
			}
		}
	}
}

func (db *DotBuilder) buildEdges() error {
	for _, e := range db.archDiagram.Edges() {
		if edge := db.buildEdge(e); edge != nil {
			db.dot.Edges = append(db.dot.Edges, edge)
		}
	}
	return nil
}

func (db *DotBuilder) buildEdge(e arch.Edge) *entity.Edge {
	fromPort := valueobject.PortStr(e.From())
	toPort := valueobject.PortStr(e.To())

	if db.isDeepMode() {
		if nodePort := db.portMap[fromPort]; nodePort != "" {
			fromPort = fmt.Sprintf("%s:%s", nodePort, fromPort)
		}
		if nodePort := db.portMap[toPort]; nodePort != "" {
			toPort = fmt.Sprintf("%s:%s", nodePort, toPort)
		}
	}
	return &entity.Edge{
		From:    fromPort,
		To:      toPort,
		Tooltip: fmt.Sprintf("%s -> %s: \n\n%s", path.Base(e.From()), path.Base(e.To()), ConcatenateRelationPos(e.Pos())),
		L:       strconv.Itoa(e.Count()),
		T:       string(dot.EdgeTypeDot),
		A:       string(db.arrowHead(e)),
	}
}

func ConcatenateRelationPos(relations []arch.RelationPos) string {
	var sb strings.Builder

	for _, relation := range relations {
		from := relation.From()
		to := relation.To()

		sb.WriteString(fmt.Sprintf("From: %s (Line: %d, Column: %d) ", path.Base(from.Filename()), from.Line(), from.Column()))

		sb.WriteString(fmt.Sprintf("To: %s (Line: %d, Column: %d)\n", path.Base(to.Filename()), to.Line(), to.Column()))
	}

	return sb.String()
}

func (db *DotBuilder) arrowHead(e arch.Edge) dot.EdgeArrowHead {
	switch e.Type() {
	case arch.RelationTypeAggregationRoot, arch.RelationTypeAggregation:
		return dot.EdgeArrowHeadDiamond
	case arch.RelationTypeAssociation:
		return dot.EdgeArrowHeadNone
	case arch.RelationTypeDependency:
		return dot.EdgeArrowHeadNormal
	}
	return dot.EdgeArrowHeadNormal
}

func (db *DotBuilder) buildTemplates() error {

	db.dot.Templates = []string{
		valueobject.TmplEdge, valueobject.TmplColumn, valueobject.TmplRow,
		valueobject.TmplSimpleNode, valueobject.TmplSimpleSubGraph, valueobject.TmplSimpleGraph,
	}

	//if db.archDiagram.NestingDepth() == 1 {
	//	db.dot.Templates = []string{
	//		valueobject.TmplEdge, valueobject.TmplSimpleNode,
	//		valueobject.TmplSimpleSubGraph, valueobject.TmplSimpleGraph,
	//	}
	//} else {
	//	db.dot.Templates = []string{
	//		valueobject.TmplColumn, valueobject.TmplRow, valueobject.TmplNode,
	//		valueobject.TmplEdge, valueobject.TmplSubGraph, valueobject.TmplGraph,
	//	}
	//}

	return nil
}

func (db *DotBuilder) isDeepMode() bool {
	return db.archDiagram.Type() == arch.TableDiagram
}
