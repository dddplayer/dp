package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"testing"
)

func TestObjMethods(t *testing.T) {
	testCases := []struct {
		name             string
		obj              *obj
		expectedName     string
		expectedPackage  string
		expectedFilename string
		expectedOffset   int
		expectedLine     int
		expectedColumn   int
	}{
		{
			name: "Object with Identifier and Position",
			obj: &obj{
				id:  &ident{name: "name1", pkg: "package1"},
				pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
			},
			expectedName:     "name1",
			expectedPackage:  "package1",
			expectedFilename: "file1.txt",
			expectedOffset:   100,
			expectedLine:     5,
			expectedColumn:   10,
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualName := tc.obj.Identifier().Name()
			actualPackage := tc.obj.Identifier().Dir()
			actualFilename := tc.obj.Position().Filename()
			actualOffset := tc.obj.Position().Offset()
			actualLine := tc.obj.Position().Line()
			actualColumn := tc.obj.Position().Column()

			if actualName != tc.expectedName || actualPackage != tc.expectedPackage ||
				actualFilename != tc.expectedFilename || actualOffset != tc.expectedOffset ||
				actualLine != tc.expectedLine || actualColumn != tc.expectedColumn {
				t.Errorf("For test case %s:\nExpected: (%s, %s, %s, %d, %d, %d)\nGot: (%s, %s, %s, %d, %d, %d)",
					tc.name, tc.expectedName, tc.expectedPackage, tc.expectedFilename, tc.expectedOffset,
					tc.expectedLine, tc.expectedColumn, actualName, actualPackage, actualFilename,
					actualOffset, actualLine, actualColumn)
			}
		})
	}
}

func TestClassMethods(t *testing.T) {
	testCases := []struct {
		name            string
		class           *Class
		attributeName   string
		methodName      string
		expectedAttrs   []string
		expectedMethods []string
	}{
		{
			name: "Class with Attributes and Methods",
			class: &Class{
				obj: &obj{
					id:  &ident{name: "Class1", pkg: "package1"},
					pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
				},
				attrs:   []*ident{{name: "attr1", pkg: "package1"}, {name: "attr2", pkg: "package1"}},
				methods: []*ident{{name: "method1", pkg: "package1"}, {name: "method2", pkg: "package1"}},
			},
			expectedAttrs:   []string{"attr1", "attr2"},
			expectedMethods: []string{"method1", "method2"},
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualAttrs := make([]string, len(tc.class.Attributes()))
			for i, attr := range tc.class.Attributes() {
				actualAttrs[i] = attr.Name()
			}
			actualMethods := make([]string, len(tc.class.Methods()))
			for i, method := range tc.class.Methods() {
				actualMethods[i] = method.Name()
			}

			if !stringSlicesEqual(actualAttrs, tc.expectedAttrs) || !stringSlicesEqual(actualMethods, tc.expectedMethods) {
				t.Errorf("For test case %s:\nExpected: (%v, %v)\nGot: (%v, %v)",
					tc.name, tc.expectedAttrs, tc.expectedMethods, actualAttrs, actualMethods)
			}
		})
	}
}

func TestInterfaceMethods(t *testing.T) {
	testCases := []struct {
		name              string
		iface             *Interface
		interfaceMethodID *ident
		expectedMethods   []string
	}{
		{
			name: "Interface with Appended Methods",
			iface: &Interface{
				obj: &obj{
					id:  &ident{name: "Interface1", pkg: "package1"},
					pos: &pos{filename: "file1.txt", offset: 100, line: 5, column: 10},
				},
			},
			interfaceMethodID: &ident{name: "method1", pkg: "package1"},
			expectedMethods:   []string{"method1"},
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interfaceMethod := &InterfaceMethod{&obj{id: tc.interfaceMethodID}}
			tc.iface.Append(interfaceMethod)

			actualMethods := make([]string, len(tc.iface.Methods()))
			for i, method := range tc.iface.Methods() {
				actualMethods[i] = method.Name()
			}

			if !stringSlicesEqual(actualMethods, tc.expectedMethods) {
				t.Errorf("For test case %s:\nExpected: %v\nGot: %v",
					tc.name, tc.expectedMethods, actualMethods)
			}
		})
	}
}

func stringSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func TestNewClass(t *testing.T) {
	mockObject := &MockObject{
		id: &ident{
			name: "MockObject",
			pkg:  "testpackage",
		},
		position: &pos{
			filename: "mockfile",
			offset:   10,
			line:     5,
			column:   2,
		},
	}

	attributeIdentifiers := []arch.ObjIdentifier{
		&ident{
			name: "attribute1",
			pkg:  "testpackage",
		},
	}

	methodIdentifiers := []arch.ObjIdentifier{
		&ident{
			name: "method1",
			pkg:  "testpackage",
		},
	}

	class := NewClass(mockObject, attributeIdentifiers, methodIdentifiers)

	if class.Identifier().Name() != "MockObject" {
		t.Errorf("Expected class name to be 'MockObject', but got '%s'", class.Identifier().Name())
	}

	if len(class.Attributes()) != len(attributeIdentifiers) {
		t.Errorf("Expected %d attributes, but got %d", len(attributeIdentifiers), len(class.Attributes()))
	}

	if len(class.Methods()) != len(methodIdentifiers) {
		t.Errorf("Expected %d methods, but got %d", len(methodIdentifiers), len(class.Methods()))
	}
}

func TestNewFunction_WithReceiver(t *testing.T) {
	// 创建一个模拟的对象标识符和位置
	mockIdentifier := &MockObject{
		id: &ident{
			name: "Object1",
			pkg:  "pkg1",
		},
		name: "Object1",
		dir:  "pkg1",
		position: &MockPosition{
			filename: "file1",
			offset:   10,
			line:     5,
			column:   2,
		},
	}

	// 调用 NewFunction 函数创建带有接收器的函数
	function := NewFunction(mockIdentifier, mockIdentifier.id)

	// 验证函数的属性是否正确设置
	if function.obj.Identifier().Name() != mockIdentifier.Name() {
		t.Errorf("Expected function.obj.Identifier().Name() to be %s, but got %s", mockIdentifier.Name(), function.obj.Identifier().Name())
	}

	if function.Receiver == nil {
		t.Error("Expected function.Receiver to be a non-nil value")
	}

	if function.Receiver.name != mockIdentifier.id.Name() {
		t.Errorf("Expected function.Receiver.name to be %s, but got %s", mockIdentifier.id.Name(), function.Receiver.name)
	}

	if function.Receiver.pkg != mockIdentifier.id.Dir() {
		t.Errorf("Expected function.Receiver.pkg to be %s, but got %s", mockIdentifier.id.Dir(), function.Receiver.pkg)
	}
}

func TestNewFunction_WithoutReceiver(t *testing.T) {
	// 创建一个模拟的对象标识符和位置
	mockIdentifier := &MockObject{
		id: &ident{
			name: "Object2",
			pkg:  "pkg2",
		},
		name: "Object2",
		dir:  "pkg2",
		position: &MockPosition{
			filename: "file2",
			offset:   20,
			line:     10,
			column:   3,
		},
	}

	// 调用 NewFunction 函数创建没有接收器的函数
	function := NewFunction(mockIdentifier, nil)

	// 验证函数的属性是否正确设置
	if function.obj.Identifier().Name() != mockIdentifier.Name() {
		t.Errorf("Expected function.obj.Identifier().Name() to be %s, but got %s", mockIdentifier.Name(), function.obj.Identifier().Name())
	}

	if function.Receiver != nil {
		t.Error("Expected function.Receiver to be nil")
	}
}

func TestNewAttr(t *testing.T) {
	// 创建一个模拟的对象
	mockObject := &MockObject{
		id: &ident{
			name: "Object1",
			pkg:  "pkg1",
		},
		name: "Object1",
		dir:  "pkg1",
		position: &MockPosition{
			filename: "file1",
			offset:   10,
			line:     5,
			column:   2,
		},
	}

	// 调用 NewAttr 函数创建属性对象
	attr := NewAttr(mockObject)

	// 验证属性对象的属性是否正确设置
	if attr.obj.Identifier().Name() != mockObject.Name() {
		t.Errorf("Expected attr.obj.Identifier().Name() to be %s, but got %s", mockObject.Name(), attr.obj.Identifier().Name())
	}
}

func TestNewInterface(t *testing.T) {
	// 创建一个模拟的对象
	mockObject := &MockObject{
		id: &ident{
			name: "Object1",
			pkg:  "pkg1",
		},
		name: "Object1",
		dir:  "pkg1",
		position: &MockPosition{
			filename: "file1",
			offset:   10,
			line:     5,
			column:   2,
		},
	}

	// 创建一个包含模拟方法的对象切片
	mockMethods := []arch.Object{
		&MockObject{
			id: &ident{
				name: "Method1",
				pkg:  "pkg1",
			},
			name: "Method1",
			dir:  "pkg1",
			position: &MockPosition{
				filename: "file2",
				offset:   20,
				line:     10,
				column:   3,
			},
		},
		// 可以添加更多的方法对象
	}

	// 调用 NewInterface 函数创建接口对象
	iface := NewInterface(mockObject, mockMethods)

	// 验证接口对象的属性是否正确设置
	if iface.obj.Identifier().Name() != mockObject.Name() {
		t.Errorf("Expected iface.obj.Identifier().Name() to be %s, but got %s", mockObject.Name(), iface.obj.Identifier().Name())
	}

	// 验证接口对象的方法是否正确设置
	if len(iface.methods) != len(mockMethods) {
		t.Errorf("Expected %d methods, but got %d", len(mockMethods), len(iface.methods))
	}
}

func TestNewGeneral(t *testing.T) {
	// 创建一个模拟的对象
	mockObject := &MockObject{
		id: &ident{
			name: "Object1",
			pkg:  "pkg1",
		},
		name: "Object1",
		dir:  "pkg1",
		position: &MockPosition{
			filename: "file1",
			offset:   10,
			line:     5,
			column:   2,
		},
	}

	// 调用 NewGeneral 函数创建 General 对象
	general := NewGeneral(mockObject)

	// 验证 General 对象的属性是否正确设置
	if general.obj.Identifier().Name() != mockObject.Name() {
		t.Errorf("Expected general.obj.Identifier().Name() to be %s, but got %s", mockObject.Name(), general.obj.Identifier().Name())
	}
}
