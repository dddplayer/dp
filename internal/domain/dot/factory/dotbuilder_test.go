package factory

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	archEntity "github.com/dddplayer/dp/internal/domain/arch/entity"
	"github.com/dddplayer/dp/internal/domain/dot"
	dotEntity "github.com/dddplayer/dp/internal/domain/dot/entity"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestNewDotBuilder(t *testing.T) {
	// 创建一个虚拟的 arch.Diagram 对象
	name := "TestDiagram"
	mockDiagram, err := archEntity.NewDiagram(name, arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 调用 NewDotBuilder 函数创建 DotBuilder 实例
	dotBuilder := NewDotBuilder(mockDiagram)

	// 验证返回的 DotBuilder 实例是否不为 nil
	if dotBuilder == nil {
		t.Error("Expected a non-nil DotBuilder instance, but got nil")
	}

	// 验证 DotBuilder 实例的 archDiagram 字段是否设置正确
	if dotBuilder.archDiagram != mockDiagram {
		t.Errorf("Expected archDiagram to be set to the mock diagram, but got %+v", dotBuilder.archDiagram)
	}

	// 验证 DotBuilder 实例的 portMap 字段是否为空
	if len(dotBuilder.portMap) != 0 {
		t.Error("Expected portMap to be an empty map, but it is not empty")
	}
}

func TestIsDeepMode(t *testing.T) {
	// 创建一个虚拟的 arch.Diagram 对象（可以是 TableDiagram 或其他类型）
	mockTableDiagram, err := archEntity.NewDiagram("MockTableDiagram", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}
	mockOtherDiagram, err := archEntity.NewDiagram("MockOtherDiagram", arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 创建 DotBuilder 实例，分别设置 archDiagram 为 TableDiagram 和其他类型
	tableModeBuilder := &DotBuilder{archDiagram: mockTableDiagram}
	otherModeBuilder := &DotBuilder{archDiagram: mockOtherDiagram}

	// 验证 isDeepMode 是否在 TableDiagram 模式下返回 true
	if !tableModeBuilder.isDeepMode() {
		t.Error("Expected isDeepMode to be true for TableDiagram, but it's false")
	}

	// 验证 isDeepMode 是否在其他模式下返回 false
	if otherModeBuilder.isDeepMode() {
		t.Error("Expected isDeepMode to be false for non-TableDiagram, but it's true")
	}
}

func TestBuildTemplates(t *testing.T) {
	// 创建 DotBuilder 实例
	dotBuilder := &DotBuilder{
		dot: &dotEntity.Dot{},
	}

	// 调用 buildTemplates 函数
	err := dotBuilder.buildTemplates()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 验证 db.dot.Templates 是否被正确设置
	expectedTemplates := []string{
		valueobject.TmplEdge, valueobject.TmplColumn, valueobject.TmplRow,
		valueobject.TmplSimpleNode, valueobject.TmplSimpleSubGraph, valueobject.TmplSimpleGraph,
	}

	// 检查 db.dot.Templates 和 expectedTemplates 是否相等
	if !reflect.DeepEqual(dotBuilder.dot.Templates, expectedTemplates) {
		t.Errorf("Expected Templates to be %v, but got %v", expectedTemplates, dotBuilder.dot.Templates)
	}
}

func TestArrowHead(t *testing.T) {
	// 创建 DotBuilder 实例
	dotBuilder := &DotBuilder{}

	// 创建不同类型的 DummyDotEdge
	aggregationRootEdge := &DummyDotEdge{FromVal: "A", ToVal: "B", T: arch.RelationTypeAggregationRoot}
	aggregationEdge := &DummyDotEdge{FromVal: "C", ToVal: "D", T: arch.RelationTypeAggregation}
	associationEdge := &DummyDotEdge{FromVal: "E", ToVal: "F", T: arch.RelationTypeAssociation}
	dependencyEdge := &DummyDotEdge{FromVal: "G", ToVal: "H", T: arch.RelationTypeDependency}
	unknownEdge := &DummyDotEdge{FromVal: "I", ToVal: "J", T: arch.RelationTypeNone}

	// 验证不同类型的 DummyDotEdge 是否返回正确的箭头头部类型
	expectedArrowHeadAggregationRoot := dot.EdgeArrowHeadDiamond
	actualArrowHeadAggregationRoot := dotBuilder.arrowHead(aggregationRootEdge)
	if actualArrowHeadAggregationRoot != expectedArrowHeadAggregationRoot {
		t.Errorf("Expected arrow head for AggregationRoot edge to be %v, but got %v", expectedArrowHeadAggregationRoot, actualArrowHeadAggregationRoot)
	}

	expectedArrowHeadAggregation := dot.EdgeArrowHeadDiamond
	actualArrowHeadAggregation := dotBuilder.arrowHead(aggregationEdge)
	if actualArrowHeadAggregation != expectedArrowHeadAggregation {
		t.Errorf("Expected arrow head for Aggregation edge to be %v, but got %v", expectedArrowHeadAggregation, actualArrowHeadAggregation)
	}

	expectedArrowHeadAssociation := dot.EdgeArrowHeadNone
	actualArrowHeadAssociation := dotBuilder.arrowHead(associationEdge)
	if actualArrowHeadAssociation != expectedArrowHeadAssociation {
		t.Errorf("Expected arrow head for Association edge to be %v, but got %v", expectedArrowHeadAssociation, actualArrowHeadAssociation)
	}

	expectedArrowHeadDependency := dot.EdgeArrowHeadNormal
	actualArrowHeadDependency := dotBuilder.arrowHead(dependencyEdge)
	if actualArrowHeadDependency != expectedArrowHeadDependency {
		t.Errorf("Expected arrow head for Dependency edge to be %v, but got %v", expectedArrowHeadDependency, actualArrowHeadDependency)
	}

	// 验证未知类型的 DummyDotEdge 是否返回默认的箭头头部类型
	expectedArrowHeadUnknown := dot.EdgeArrowHeadNormal
	actualArrowHeadUnknown := dotBuilder.arrowHead(unknownEdge)
	if actualArrowHeadUnknown != expectedArrowHeadUnknown {
		t.Errorf("Expected arrow head for Unknown edge to be %v, but got %v", expectedArrowHeadUnknown, actualArrowHeadUnknown)
	}
}

func TestConcatenateRelationPos(t *testing.T) {
	// 创建一些模拟的 RelationPos 对象
	relPos1 := &MockRelationPos{
		fromPos: &MockPosition{filename: "/path/to/file1.txt", line: 10, column: 5},
		toPos:   &MockPosition{filename: "/path/to/file2.txt", line: 20, column: 15}}
	relPos2 := &MockRelationPos{
		fromPos: &MockPosition{filename: "/path/to/file3.txt", line: 30, column: 25},
		toPos:   &MockPosition{filename: "/path/to/file4.txt", line: 40, column: 35}}

	// 将 RelationPos 对象放入切片中
	relations := []arch.RelationPos{relPos1, relPos2}

	// 调用 ConcatenateRelationPos 函数
	result := ConcatenateRelationPos(relations)

	// 验证结果是否符合预期
	expectedResult := "From: file1.txt (Line: 10, Column: 5) To: file2.txt (Line: 20, Column: 15)\n" +
		"From: file3.txt (Line: 30, Column: 25) To: file4.txt (Line: 40, Column: 35)\n"

	if result != expectedResult {
		t.Errorf("Expected result:\n%s\nBut got:\n%s", expectedResult, result)
	}
}

func TestDotBuilder_buildEdge(t *testing.T) {
	// 创建一个虚拟的 DummyDotEdge 对象
	dummyEdge := &DummyDotEdge{
		FromVal: "NodeA",
		ToVal:   "NodeB",
		L:       "2",
		TT:      "Dot",
		T:       arch.RelationTypeDependency, // 假设关系类型为 Dependency
	}

	mockDiagram, err := archEntity.NewDiagram("MockOtherDiagram", arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 创建一个虚拟的 DotBuilder 实例
	dotBuilder := &DotBuilder{
		archDiagram: mockDiagram, // 假设已经有一个模拟的 arch.Diagram 对象
		portMap:     make(map[string]string),
	}

	// 调用 buildEdge 函数生成 entity.Edge
	edge := dotBuilder.buildEdge(dummyEdge)

	// 验证生成的 entity.Edge 是否符合预期
	expectedEdge := &dotEntity.Edge{
		From:    valueobject.PortStr(dummyEdge.FromVal),
		To:      valueobject.PortStr(dummyEdge.ToVal),
		Tooltip: fmt.Sprintf("%s -> %s: \n\n%s", path.Base(dummyEdge.FromVal), path.Base(dummyEdge.ToVal), ConcatenateRelationPos([]arch.RelationPos{})), // 假设传入空的关系位置
		L:       "1",
		T:       "solid",
		A:       string(dot.EdgeArrowHeadNormal), // 假设关系类型为 Dependency
	}

	if !reflect.DeepEqual(edge, expectedEdge) {
		t.Errorf("Generated edge does not match the expected result.\nExpected: %+v\nGot: %+v", expectedEdge, edge)
	}
}

func TestDotBuilder_buildEdges(t *testing.T) {
	// 创建一个虚拟的 arch.Diagram 对象
	mockDiagram, err := archEntity.NewDiagram("test", arch.TableDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if err := mockDiagram.AddStringTo("mockGroup", mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 创建一个虚拟的 DotBuilder 实例
	dotBuilder := &DotBuilder{
		archDiagram: mockDiagram,
		portMap:     make(map[string]string),
		dot:         &dotEntity.Dot{},
	}

	// 调用 buildEdges 函数
	if err := dotBuilder.buildEdges(); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 验证生成的 entity.Edge 是否符合预期
	expectedEdges := []dotEntity.Edge{
		// 与模拟的边缘匹配的期望 entity.Edge 对象
		{
			From: "test",
			To:   "mockGroup",
		},
		// 添加更多期望的 entity.Edge 对象
	}

	// 验证生成的 entity.Edge 是否与预期的结果匹配
	if len(mockDiagram.Edges()) != len(expectedEdges) {
		t.Errorf("Generated edges do not match the expected results.\nExpected: %+v\nGot: %+v", expectedEdges, dotBuilder.dot.Edges)
	}
	if mockDiagram.Edges()[0].From() != expectedEdges[0].From {
		t.Errorf("Generated edges do not match the expected results.\nExpected: %+v\nGot: %+v", expectedEdges, dotBuilder.dot.Edges)
	}
	if mockDiagram.Edges()[0].To() != expectedEdges[0].To {
		t.Errorf("Generated edges do not match the expected results.\nExpected: %+v\nGot: %+v", expectedEdges, dotBuilder.dot.Edges)
	}
}

func TestDotBuilder_buildPortMap(t *testing.T) {
	// 创建一个虚拟的 entity.SubGraph 对象
	mockSubGraph := &dotEntity.SubGraph{
		// 初始化模拟的节点和表格数据
		Nodes: []*dotEntity.Node{
			{
				ID: "NodeA",
				Table: &dotEntity.Table{
					Rows: []*dotEntity.Row{
						{
							Data: []*dotEntity.Data{
								{Port: "Port1"},
								{Port: "Port2"},
								// 添加更多模拟的数据
							},
						},
					},
				},
			},
			// 添加更多模拟的节点
		},
	}

	// 创建一个虚拟的 DotBuilder 实例
	dotBuilder := &DotBuilder{
		archDiagram: nil, // 不需要关注这个字段
		portMap:     make(map[string]string),
		dot:         &dotEntity.Dot{},
	}

	// 调用 buildPortMap 函数
	dotBuilder.buildPortMap(mockSubGraph)

	// 验证生成的 portMap 是否符合预期
	expectedPortMap := map[string]string{
		"Port1": "NodeA",
		"Port2": "NodeA",
		// 添加更多期望的映射
	}

	// 验证生成的 portMap 是否与预期的结果匹配
	if !reflect.DeepEqual(dotBuilder.portMap, expectedPortMap) {
		t.Errorf("Generated portMap does not match the expected results.\nExpected: %+v\nGot: %+v", expectedPortMap, dotBuilder.portMap)
	}
}

func TestDotBuilder_buildGraph(t *testing.T) {
	// 创建一个虚拟的 arch.SubDiagram 对象
	mockSubDiagram := generateDummyDotGraph() // 根据您的数据结构创建虚拟对象

	mockDiagram, err := archEntity.NewDiagram("test", arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// 创建一个虚拟的 DotBuilder 实例
	dotBuilder := &DotBuilder{
		archDiagram: mockDiagram,
		portMap:     make(map[string]string),
		dot:         &dotEntity.Dot{},
	}

	// 创建一个虚拟的 entity.SubGraph 作为父子图
	parentSubGraph := &dotEntity.SubGraph{
		// 初始化父子图的属性
		Name:      "ParentSubGraph",
		Label:     "ParentSubGraph",
		Nodes:     []*dotEntity.Node{},
		SubGraphs: []*dotEntity.SubGraph{},
	}

	// 调用 buildGraph 函数
	if err := dotBuilder.buildGraph(mockSubDiagram, parentSubGraph); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 验证生成的子图是否符合预期
	expectedSubGraph := &dotEntity.SubGraph{
		// 初始化预期的子图属性
		Name:  "dcPXfIp",   // 根据您的数据结构设置名称
		Label: "TestGraph", // 根据您的数据结构设置标签
		Nodes: []*dotEntity.Node{
			{ID: "dc0Zpf6", Name: "Node1", BgColor: "", Table: nil},
			{ID: "dc4oA8n", Name: "Node2", BgColor: "", Table: nil},
		},
		SubGraphs: []*dotEntity.SubGraph{
			// 根据您的数据结构设置预期的子图
		},
	}

	// 验证生成的子图是否与预期的结果匹配
	if !reflect.DeepEqual(parentSubGraph.SubGraphs[0], expectedSubGraph) {
		t.Errorf("Generated subgraph does not match the expected results.\nExpected: %+v\nGot: %+v", expectedSubGraph, parentSubGraph.SubGraphs[0])
	}
}

func TestDotBuilder_buildSubGraph(t *testing.T) {
	mockDiagram, err := archEntity.NewDiagram("TestDiagram", arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}
	if err := mockDiagram.AddStringTo("mockGroup", mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 在 DotBuilder 实例中设置 isDeepMode 为 true
	dotBuilder := &DotBuilder{
		archDiagram: mockDiagram,
		portMap:     make(map[string]string),
		dot:         &dotEntity.Dot{},
	}

	// 调用 buildSubGraph 函数
	if err := dotBuilder.buildSubGraph(); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 验证生成的子图是否符合预期
	expectedSubGraph1 := &dotEntity.SubGraph{
		Name:      "dbWlszi",     // 请根据您的数据结构设置名称
		Label:     "TestDiagram", // 请根据您的数据结构设置标签
		Nodes:     []*dotEntity.Node{{}},
		SubGraphs: []*dotEntity.SubGraph{{}},
	}

	sg := dotBuilder.dot.SubGraphs[0]
	if sg.Name != expectedSubGraph1.Name {
		t.Errorf("Generated subgraph does not match the expected results.\nExpected: %+v\nGot: %+v", expectedSubGraph1, sg)
	}
	if sg.Label != expectedSubGraph1.Label {
		t.Errorf("Generated subgraph does not match the expected results.\nExpected: %+v\nGot: %+v", expectedSubGraph1, sg)
	}
	if len(sg.SubGraphs) != len(expectedSubGraph1.SubGraphs) {
		t.Errorf("Generated subgraph does not match the expected results.\nExpected: %+v\nGot: %+v", expectedSubGraph1, sg)
	}
	if len(sg.Nodes) != len(expectedSubGraph1.Nodes) {
		t.Errorf("Generated subgraph does not match the expected results.\nExpected: %+v\nGot: %+v", expectedSubGraph1, sg)
	}
}

func TestDotBuilder_Build(t *testing.T) {
	// 创建虚拟的 arch.Diagram 对象
	mockDiagram, err := archEntity.NewDiagram("TestDiagram", arch.PlainDiagram)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}
	if err := mockDiagram.AddStringTo("mockGroup", mockDiagram.Name(), arch.RelationTypeAggregationRoot); err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 创建一个虚拟的 DotBuilder 实例
	dotBuilder := &DotBuilder{
		archDiagram: mockDiagram,
		portMap:     make(map[string]string),
		dot:         &dotEntity.Dot{},
	}

	// 调用 Build 函数
	d, err := dotBuilder.Build()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}

	// 验证生成的 Dot 实例是否符合预期
	expectedDot := &dotEntity.Dot{
		Name:      "TestDiagram",                                            // 根据您的数据结构设置名称
		Label:     "\n\nTestDiagram\nDomain Model\n\nPowered by DDD Player", // 根据您的数据结构设置标签
		SubGraphs: []*dotEntity.SubGraph{{}},                                // 根据您的数据结构设置子图
		Edges:     []*dotEntity.Edge{{}},                                    // 根据您的数据结构设置边
	}

	if d.Name != expectedDot.Name {
		t.Errorf("Generated dot does not match the expected results.\nExpected: %+v\nGot: %+v", expectedDot, d)
	}
	if !strings.HasPrefix(d.Label, expectedDot.Label) {
		t.Errorf("Generated dot does not match the expected results.\nExpected: %+v\nGot: %+v", expectedDot, d)
	}
}
