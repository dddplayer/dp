package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/code"
	"reflect"
	"testing"
)

func TestDomainModel_HandleClass(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifier and position
	id := &ident{name: "MyClass", pkg: "domain/myclass"}
	pos := &pos{filename: "myclass.go", offset: 10, line: 5, column: 15}

	// call handleClass method
	dm.handleClass(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj arch.Object
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*Class); !ok {
		t.Errorf("Expected object in repository to be a Class, but got %T", obj)
	}

	class := obj.(*Class)
	if class.obj.id.name != id.name || class.obj.id.pkg != id.pkg {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, class.obj.id)
	}

	if class.obj.pos.filename != pos.filename ||
		class.obj.pos.offset != pos.offset || class.obj.pos.line != pos.line ||
		class.obj.pos.column != pos.column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, class.obj.pos)
	}
}

func TestDomainModel_HandleGenObj(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifier and position
	id := &ident{name: "MyObj", pkg: "domain/myobj"}
	pos := &pos{filename: "myobj.go", offset: 10, line: 5, column: 15}

	// call handleGenObj method
	dm.handleGenObj(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj arch.Object
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*General); !ok {
		t.Errorf("Expected object in repository to be a General, but got %T", obj)
	}

	genObj := obj.(*General)
	if genObj.obj.id.name != id.name || genObj.obj.id.pkg != id.pkg {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, genObj.obj.id)
	}

	if genObj.obj.pos.filename != pos.filename ||
		genObj.obj.pos.offset != pos.offset || genObj.obj.pos.line != pos.line ||
		genObj.obj.pos.column != pos.column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, genObj.obj.pos)
	}
}

func TestDomainModel_HandleFunc(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifier and position
	id := &ident{name: "myFunction", pkg: "domain/myFunction"}
	pos := &pos{filename: "myFunction.go", offset: 10, line: 5, column: 15}

	// call handleFunc method
	dm.handleFunc(id, pos, nil)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj arch.Object
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*Function); !ok {
		t.Errorf("Expected object in repository to be a Function, but got %T", obj)
	}

	fn := obj.(*Function)
	if fn.obj.id.name != id.name || fn.obj.id.pkg != id.pkg {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, fn.obj.id)
	}

	if fn.obj.pos.filename != pos.filename ||
		fn.obj.pos.offset != pos.offset || fn.obj.pos.line != pos.line ||
		fn.obj.pos.column != pos.column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, fn.obj.pos)
	}
}

func TestDomainModel_LinkHandler(t *testing.T) {
	// mock Repository
	repo := newMockRelationRepository()

	// create Handler
	dm := &CodeHandler{Scope: "test", RelRepo: repo}

	// create link
	fromId := &ident{name: "field", pkg: "/test/myField"}
	fromPos := &pos{filename: "class.go", offset: 10, line: 5, column: 15}
	toId := &ident{name: "func", pkg: "/test/myFunction"}

	link := &code.Link{
		From: &code.Node{
			Meta:   newDummyMetaWithIdent(fromId),
			Pos:    fromPos,
			Type:   code.TypeGenStructField,
			Parent: nil,
		},
		To: &code.Node{
			Meta: newDummyMetaWithIdent(toId),
			Type: code.TypeAny,
		},
		Relation: code.OneOne,
	}

	// call LinkHandler method
	dm.LinkHandler(link)

	var rel arch.Relation
	if rel = repo.Find(fromId); rel == nil {
		t.Errorf("Expected repository to have relation %v, but got %v", fromId, rel)
	}

	if obj, ok := rel.(*Association); !ok {
		t.Errorf("Expected object in repository to be a Association, but got %T", obj)
	}

	relation := rel.(*Association)
	if relation.From().Identifier().Name() != fromId.name || relation.Refer().Identifier().Name() != toId.name {
		t.Errorf("Expected object in repository to have From %v and To %v, but got From %v and To %v",
			fromId, toId, relation.From().Identifier().Name(), relation.Refer().Identifier().Name())
	}

	if relation.From().Position().Filename() != fromPos.Filename() ||
		relation.From().Position().Offset() != fromPos.Offset() ||
		relation.From().Position().Line() != fromPos.Line() ||
		relation.From().Position().Column() != fromPos.Column() {
		t.Errorf("Expected object in repository to have pos %v, but got %v", fromPos, relation.From().Position())
	}

	if relation.AssociationType() != arch.RelationType(link.Relation) {
		t.Errorf("Expected object in repository to have ship %v, but got %v", link.Relation, relation.AssociationType())
	}
}

func TestDomainModel_NodeHandler(t *testing.T) {
	// Create a new mock repository
	repo := newMockRepository()

	// Create a new domain model with the mock repository
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create link
	id1 := &ident{name: "id1", pkg: "/test/id1"}
	pos1 := &pos{filename: "class.go", offset: 10, line: 5, column: 15}
	cid1 := &ident{name: id1.name, pkg: id1.pkg}
	cpos1 := &pos{filename: pos1.filename, offset: pos1.offset, line: pos1.line, column: pos1.column}
	idEntity := &ident{name: "entity1", pkg: "/ddd/entity"}
	cidEntity := &ident{name: idEntity.name, pkg: idEntity.pkg}
	idValueObject := &ident{name: "valueobject1", pkg: "/ddd/valueobject"}
	cidValueObject := &ident{name: idValueObject.name, pkg: idValueObject.pkg}
	idInterface := &ident{name: "interface1", pkg: "/ddd/interface"}
	cidInterface := &ident{name: idInterface.name, pkg: idInterface.pkg}

	// Define the test cases
	tests := []struct {
		name     string
		node     *code.Node
		want     arch.Object
		wantType reflect.Type
	}{
		{
			name:     "Test handleGenObj for TypeGenIdent",
			node:     &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeGenIdent},
			want:     &General{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&General{}),
		},
		{
			name:     "Test handleGenObj for TypeGenIdent",
			node:     &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeGenFunc},
			want:     &General{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&General{}),
		},
		{
			name:     "Test handle entity for TypeGenStruct",
			node:     &code.Node{Meta: newDummyMetaWithIdent(idEntity), Pos: pos1, Type: code.TypeGenStruct},
			want:     &Class{obj: &obj{id: cidEntity, pos: cpos1}},
			wantType: reflect.TypeOf(&Class{}),
		},
		{
			name:     "Test handle entity for TypeGenStruct",
			node:     &code.Node{Meta: newDummyMetaWithIdent(idValueObject), Pos: pos1, Type: code.TypeGenStruct},
			want:     &Class{obj: &obj{id: cidValueObject, pos: cpos1}},
			wantType: reflect.TypeOf(&Class{}),
		},
		{
			name: "Test handle struct field for TypeGenStructField",
			node: &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeGenStructField,
				Parent: &code.Node{Meta: newDummyMetaWithIdent(idValueObject), Pos: pos1, Type: code.TypeGenStruct}},
			want:     &Attr{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&Attr{}),
		},
		{
			name:     "Test handle interface for TypeGenInterface",
			node:     &code.Node{Meta: newDummyMetaWithIdent(idInterface), Pos: pos1, Type: code.TypeGenInterface},
			want:     &Interface{obj: &obj{id: cidInterface, pos: cpos1}, methods: nil},
			wantType: reflect.TypeOf(&Interface{}),
		},
		{
			name: "Test handle interface method for TypeGenInterfaceMethod",
			node: &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeGenInterfaceMethod,
				Parent: &code.Node{Meta: newDummyMetaWithIdent(idInterface), Pos: pos1, Type: code.TypeGenInterface}},
			want:     &InterfaceMethod{&obj{id: id1, pos: cpos1}},
			wantType: reflect.TypeOf(&InterfaceMethod{}),
		},
		{
			name: "Test handle function with parent non nil for TypeFunc",
			node: &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeFunc,
				Parent: &code.Node{Meta: newDummyMetaWithIdent(idEntity), Pos: pos1, Type: code.TypeGenStruct}},
			want:     &Function{obj: &obj{id: id1, pos: cpos1}},
			wantType: reflect.TypeOf(&Function{}),
		},
		{
			name:     "Test handle function with parent nil for TypeFunc",
			node:     &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeFunc, Parent: nil},
			want:     &Function{obj: &obj{id: id1, pos: cpos1}},
			wantType: reflect.TypeOf(&Function{}),
		},
		{
			name:     "Test handle other type for default",
			node:     &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeNone, Parent: nil},
			want:     &General{&obj{id: cid1, pos: cpos1}},
			wantType: reflect.TypeOf(&General{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.Clear()
			dm.NodeHandler(tt.node)

			// Check if the object was correctly inserted into the repository
			got := repo.Find(&ident{
				name: tt.node.Meta.Name(),
				pkg:  tt.node.Meta.Pkg(),
			})

			if reflect.TypeOf(got) != tt.wantType {
				t.Errorf("NodeHandler() = type %v, want %v", reflect.TypeOf(got), tt.wantType)
			}

			if got.Identifier().ID() != tt.want.Identifier().ID() {
				t.Errorf("NodeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
