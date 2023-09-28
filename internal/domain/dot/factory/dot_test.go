package factory

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/dot/entity"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"path"
	"reflect"
	"testing"
)

func TestNewSubGraph(t *testing.T) {
	dummySubDiagram := generateDummyDotGraph()

	testCases := []struct {
		name             string
		subDiagram       arch.SubDiagram
		expectedSubGraph *entity.SubGraph
	}{
		{
			name:       "SubDiagram with Nodes",
			subDiagram: dummySubDiagram,
			expectedSubGraph: &entity.SubGraph{
				Name:  valueobject.GenerateShortURL("TestGraph"),
				Label: "TestGraph",
				Nodes: []*entity.Node{
					{
						ID:      valueobject.GenerateShortURL("Node1"),
						Name:    "Node1",
						BgColor: "",
						Table:   nil,
					},
					{
						ID:      valueobject.GenerateShortURL("Node2"),
						Name:    "Node2",
						BgColor: "",
						Table:   nil,
					},
				},
				SubGraphs: []*entity.SubGraph{},
			},
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualSubGraph := NewSubGraph(tc.subDiagram)
			if !reflect.DeepEqual(actualSubGraph, tc.expectedSubGraph) {
				t.Errorf("For test case %s:\nExpected: %+v\nGot: %+v",
					tc.name, tc.expectedSubGraph, actualSubGraph)
			}
		})
	}
}

func TestNewSummarySubGraph(t *testing.T) {
	// 创建虚拟的 arch.SubDiagram 对象（包含摘要信息的和不包含摘要信息的）
	dummySubDiagram := generateDummyDotGraph()

	testCases := []struct {
		name             string
		subDiagram       arch.SubDiagram
		expectedSubGraph *entity.SubGraph
	}{
		{
			name:       "SubDiagram without Summary",
			subDiagram: dummySubDiagram,
			expectedSubGraph: &entity.SubGraph{
				Name:  valueobject.PortStr("TestGraph"),
				Label: "TestGraph",
				Nodes: []*entity.Node{
					{
						ID:    valueobject.PortStr(dummySubDiagram.Name()),
						Name:  path.Base(dummySubDiagram.Name()),
						Table: &entity.Table{Rows: []*entity.Row{{}, {}}},
					},
				},
				SubGraphs: []*entity.SubGraph{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualSubGraph := NewSummarySubGraph(tc.subDiagram)
			if actualSubGraph.Name != tc.expectedSubGraph.Name ||
				actualSubGraph.Label != tc.expectedSubGraph.Label ||
				len(actualSubGraph.Nodes) != len(tc.expectedSubGraph.Nodes) {
				t.Errorf("For test case %s:\nExpected: %+v\nGot: %+v",
					tc.name, tc.expectedSubGraph, actualSubGraph)
			}
		})
	}
}
