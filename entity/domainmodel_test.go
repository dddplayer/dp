package entity

import (
	ca "github.com/dddplayer/core/codeanalysis/entity"
	"github.com/dddplayer/core/valueobject"
	"reflect"
	"testing"
)

type mockRepository struct {
	data map[valueobject.Identifier]DomainObject
}

func (r *mockRepository) Find(id valueobject.Identifier) DomainObject {
	return r.data[id]
}

func (r *mockRepository) Insert(obj DomainObject) error {
	r.data[*obj.Identifier()] = obj
	return nil
}

func (r *mockRepository) Walk(cb func(obj DomainObject) error) {
	for _, obj := range r.data {
		if err := cb(obj); err != nil {
			return
		}
	}
}
func TestDomainModel_HandleClass(t *testing.T) {
	// mock Repository
	repo := &mockRepository{data: map[valueobject.Identifier]DomainObject{}}

	// create DomainModel
	dm := &DomainModel{Name: "Test Model", Repo: repo}

	// create identifier and position
	id := valueobject.Identifier{Name: "MyClass", Path: "domain/myclass"}
	pos := valueobject.Position{Filename: "myclass.go", Offset: 10, Line: 5, Column: 15}

	// call handleClass method
	dm.handleClass(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj DomainObject
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*Class); !ok {
		t.Errorf("Expected object in repository to be a Class, but got %T", obj)
	}

	class := obj.(*Class)
	if class.obj.id.Name != id.Name || class.obj.id.Path != id.Path {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, class.obj.id)
	}

	if class.obj.pos.Filename != pos.Filename ||
		class.obj.pos.Offset != pos.Offset || class.obj.pos.Line != pos.Line ||
		class.obj.pos.Column != pos.Column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, class.obj.pos)
	}
}

func TestDomainModel_HandleGenObj(t *testing.T) {
	// mock Repository
	repo := &mockRepository{data: map[valueobject.Identifier]DomainObject{}}

	// create DomainModel
	dm := &DomainModel{Name: "Test Model", Repo: repo}

	// create identifier and position
	id := valueobject.Identifier{Name: "MyObj", Path: "domain/myobj"}
	pos := valueobject.Position{Filename: "myobj.go", Offset: 10, Line: 5, Column: 15}

	// call handleGenObj method
	dm.handleGenObj(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj DomainObject
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*General); !ok {
		t.Errorf("Expected object in repository to be a General, but got %T", obj)
	}

	genObj := obj.(*General)
	if genObj.obj.id.Name != id.Name || genObj.obj.id.Path != id.Path {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, genObj.obj.id)
	}

	if genObj.obj.pos.Filename != pos.Filename ||
		genObj.obj.pos.Offset != pos.Offset || genObj.obj.pos.Line != pos.Line ||
		genObj.obj.pos.Column != pos.Column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, genObj.obj.pos)
	}
}

func TestDomainModel_HandleFunc(t *testing.T) {
	// mock Repository
	repo := &mockRepository{data: map[valueobject.Identifier]DomainObject{}}

	// create DomainModel
	dm := &DomainModel{Name: "Test Model", Repo: repo}

	// create identifier and position
	id := valueobject.Identifier{Name: "myFunction", Path: "domain/myFunction"}
	pos := valueobject.Position{Filename: "myFunction.go", Offset: 10, Line: 5, Column: 15}

	// call handleFunc method
	dm.handleFunc(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj DomainObject
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*Function); !ok {
		t.Errorf("Expected object in repository to be a Function, but got %T", obj)
	}

	fn := obj.(*Function)
	if fn.obj.id.Name != id.Name || fn.obj.id.Path != id.Path {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, fn.obj.id)
	}

	if fn.obj.pos.Filename != pos.Filename ||
		fn.obj.pos.Offset != pos.Offset || fn.obj.pos.Line != pos.Line ||
		fn.obj.pos.Column != pos.Column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, fn.obj.pos)
	}
}

func TestDomainModel_LinkHandler(t *testing.T) {
	// mock Repository
	repo := &mockRepository{data: map[valueobject.Identifier]DomainObject{}}

	// create DomainModel
	dm := &DomainModel{Name: "test", Repo: repo}

	// create link
	fromId := &MockIdentifier{mockName: "field", mockPath: "/test/myField"}
	fromPos := &MockPosition{filename: "class.go", offset: 10, line: 5, column: 15}
	toId := &MockIdentifier{mockName: "func", mockPath: "/test/myFunction"}

	link := &ca.Link{
		From: &ca.Node{
			ID:     fromId,
			Pos:    fromPos,
			Type:   ca.TypeGenStructField,
			Parent: nil,
		},
		To: &ca.Node{
			ID:   toId,
			Type: ca.TypeFunc,
		},
		Relation: ca.OneOne,
	}

	// call LinkHandler method
	dm.LinkHandler(link)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj DomainObject
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*valueobject.Relation); !ok {
		t.Errorf("Expected object in repository to be a Relation, but got %T", obj)
	}

	relation := obj.(*valueobject.Relation)
	if relation.From.Name != fromId.mockName || relation.To.Name != toId.mockName {
		t.Errorf("Expected object in repository to have From %v and To %v, but got From %v and To %v", fromId, toId, relation.From.Name, relation.To.Name)
	}

	if relation.Pos.Filename != fromPos.Filename() ||
		relation.Pos.Offset != fromPos.Offset() || relation.Pos.Line != fromPos.Line() ||
		relation.Pos.Column != fromPos.Column() {
		t.Errorf("Expected object in repository to have pos %v, but got %v", fromPos, relation.Pos)
	}

	if relation.Ship != valueobject.RelationShip(link.Relation) {
		t.Errorf("Expected object in repository to have ship %v, but got %v", link.Relation, relation.Ship)
	}

	switch link.From.Type {
	case ca.TypeGenStructField:
		if relation.Type != valueobject.TypeRefer {
			t.Errorf("Expected object in repository to have TypeRefer, but got %v", relation.Type)
		}
	case ca.TypeGenInterface:
		if relation.Type != valueobject.TypeImplements {
			t.Errorf("Expected object in repository to have TypeImplements, but got %v", relation.Type)
		}
	case ca.TypeFunc:
		if relation.Type != valueobject.TypeCall {
			t.Errorf("Expected object in repository to have TypeCall, but got %v", relation.Type)
		}
	}
}

func TestDomainModel_NodeHandler(t *testing.T) {
	// Create a new mock repository
	repo := &mockRepository{data: make(map[valueobject.Identifier]DomainObject)}

	// Create a new domain model with the mock repository
	dm := &DomainModel{Name: "Test Model", Repo: repo}

	// create link
	id1 := &MockIdentifier{mockName: "id1", mockPath: "/test/id1"}
	pos1 := &MockPosition{filename: "class.go", offset: 10, line: 5, column: 15}
	cid1 := &valueobject.Identifier{Name: id1.mockName, Path: id1.mockPath}
	cpos1 := &valueobject.Position{Filename: pos1.filename, Offset: pos1.offset, Line: pos1.line, Column: pos1.column}
	idEntity := &MockIdentifier{mockName: "entity1", mockPath: "/ddd/entity"}
	cidEntity := &valueobject.Identifier{Name: idEntity.mockName, Path: idEntity.mockPath}
	idValueObject := &MockIdentifier{mockName: "valueobject1", mockPath: "/ddd/valueobject"}
	cidValueObject := &valueobject.Identifier{Name: idValueObject.mockName, Path: idValueObject.mockPath}
	idInterface := &MockIdentifier{mockName: "interface1", mockPath: "/ddd/interface"}
	cidInterface := &valueobject.Identifier{Name: idInterface.mockName, Path: idInterface.mockPath}
	idFactory := &MockIdentifier{mockName: "factory1", mockPath: "/ddd/factory"}
	cidFactory := &valueobject.Identifier{Name: idFactory.mockName, Path: idFactory.mockPath}
	idService := &MockIdentifier{mockName: "service1", mockPath: "/ddd/service"}
	cidService := &valueobject.Identifier{Name: idService.mockName, Path: idService.mockPath}

	// Define the test cases
	tests := []struct {
		name     string
		node     *ca.Node
		want     DomainObject
		wantType reflect.Type
	}{
		{
			name:     "Test handleGenObj for TypeGenIdent",
			node:     &ca.Node{ID: id1, Pos: pos1, Type: ca.TypeGenIdent},
			want:     &General{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&General{}),
		},
		{
			name:     "Test handleGenObj for TypeGenIdent",
			node:     &ca.Node{ID: id1, Pos: pos1, Type: ca.TypeGenFunc},
			want:     &General{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&General{}),
		},
		{
			name:     "Test handle entity for TypeGenStruct",
			node:     &ca.Node{ID: idEntity, Pos: pos1, Type: ca.TypeGenStruct},
			want:     &General{&obj{id: cidEntity, pos: cpos1}},
			wantType: reflect.TypeOf(&Entity{}),
		},
		{
			name:     "Test handle entity for TypeGenStruct",
			node:     &ca.Node{ID: idValueObject, Pos: pos1, Type: ca.TypeGenStruct},
			want:     &General{&obj{id: cidValueObject, pos: cpos1}},
			wantType: reflect.TypeOf(&ValueObject{}),
		},
		{
			name: "Test handle struct field for TypeGenStructField",
			node: &ca.Node{ID: id1, Pos: pos1, Type: ca.TypeGenStructField,
				Parent: &ca.Node{ID: idValueObject, Pos: pos1, Type: ca.TypeGenStruct}},
			want:     &Attr{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&Attr{}),
		},
		{
			name:     "Test handle interface for TypeGenInterface",
			node:     &ca.Node{ID: idInterface, Pos: pos1, Type: ca.TypeGenInterface},
			want:     &Interface{obj: &obj{id: cidInterface, pos: cpos1}, Methods: nil},
			wantType: reflect.TypeOf(&Interface{}),
		},
		{
			name: "Test handle interface method for TypeGenInterfaceMethod",
			node: &ca.Node{ID: id1, Pos: pos1, Type: ca.TypeGenInterfaceMethod,
				Parent: &ca.Node{ID: idInterface, Pos: pos1, Type: ca.TypeGenInterface}},
			want:     &InterfaceMethod{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&InterfaceMethod{}),
		},
		{
			name: "Test handle function with parent non nil for TypeFunc",
			node: &ca.Node{ID: id1, Pos: pos1, Type: ca.TypeFunc,
				Parent: &ca.Node{ID: idEntity, Pos: pos1, Type: ca.TypeGenStruct}},
			want:     &Function{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&Function{}),
		},
		{
			name:     "Test handle function with parent nil for TypeFunc",
			node:     &ca.Node{ID: id1, Pos: pos1, Type: ca.TypeFunc, Parent: nil},
			want:     &Function{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&Function{}),
		},
		{
			name:     "Test handle factory for TypeFunc",
			node:     &ca.Node{ID: idFactory, Pos: pos1, Type: ca.TypeFunc, Parent: nil},
			want:     &Factory{Function{&obj{id: cidFactory, pos: cpos1}}},
			wantType: reflect.TypeOf(&Factory{}),
		},
		{
			name:     "Test handle service for TypeFunc",
			node:     &ca.Node{ID: idService, Pos: pos1, Type: ca.TypeFunc, Parent: nil},
			want:     &Service{Function{&obj{id: cidService, pos: cpos1}}},
			wantType: reflect.TypeOf(&Service{}),
		},
		{
			name:     "Test handle other type for default",
			node:     &ca.Node{ID: id1, Pos: pos1, Type: ca.NodeType("non-exist"), Parent: nil},
			want:     &General{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&General{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm.NodeHandler(tt.node)

			// Check if the object was correctly inserted into the repository
			got := repo.Find(valueobject.Identifier{
				Name: tt.node.ID.Name(),
				Path: tt.node.ID.Path(),
			})

			if reflect.TypeOf(got) != tt.wantType {
				t.Errorf("NodeHandler() = type %v, want %v", reflect.TypeOf(got), tt.wantType)
			}

			if got.Identifier().String() != tt.want.Identifier().String() {
				t.Errorf("NodeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
