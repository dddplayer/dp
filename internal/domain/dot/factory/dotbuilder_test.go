package factory

//
//func TestGraphBuilder_BuildEdge(t *testing.T) {
//	// 设置测试用例
//	testCases := []struct {
//		name         string
//		dotEdge      arch.Edge
//		expectedEdge *entity.Edge
//	}{
//		{
//			name: "Test case 1",
//			dotEdge: &DummyDotEdge{
//				FromVal: "node1",
//				ToVal:   "node2",
//			},
//			expectedEdge: &entity.Edge{
//				From: "node1",
//				To:   "node2",
//				T:    "solid",
//				A:    "onormal",
//			},
//		},
//		{
//			name: "Test case 2",
//			dotEdge: &DummyDotEdge{
//				FromVal: "node3",
//				ToVal:   "node4",
//			},
//			expectedEdge: &entity.Edge{
//				From: "node3",
//				To:   "node4",
//				T:    "solid",
//				A:    "onormal",
//			},
//		},
//		// 添加更多的测试用例...
//	}
//
//	// 执行测试用例
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			// 初始化 DotBuilder
//			gb := &DotBuilder{
//				dotGraph: &DummyDotGraph{},
//			}
//
//			// 调用 buildEdge 函数，得到实际的输出
//			actualEdge := gb.buildEdge(tc.dotEdge)
//
//			// 对比实际输出和期望输出
//			if !reflect.DeepEqual(actualEdge, tc.expectedEdge) {
//				t.Errorf("Unexpected result:\nexpected: %v\nactual: %v\n", tc.expectedEdge, actualEdge)
//			}
//		})
//	}
//}

//func TestGraphBuilder_BuildNode(t *testing.T) {
//	graph := generateDummyDotGraph()
//	gb := DotBuilder{dotGraph: graph}
//
//	node1 := graph.NodesVal[0]
//	got := gb.buildNode(node1)
//
//	// Test if node is not nil
//	if got == nil {
//		t.Error("Expected a non-nil node, but got nil")
//	}
//
//	// Test if node name is correct
//	if got.Name != node1.Name() {
//		t.Errorf("Expected node name %q, but got %q", node1.Name(), got.Name)
//	}
//}
