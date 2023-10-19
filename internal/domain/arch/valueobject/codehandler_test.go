package valueobject

import (
	"errors"
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

	// call handleClass method
	dm.handleClass(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
	}
}

func TestDomainModel_HandleAttribute(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifiers and positions
	id := &ident{name: "MyAttr", pkg: "domain/myattr"}
	pos := &pos{filename: "myattr.go", offset: 10, line: 5, column: 15}
	pid := &ident{name: "MyClass", pkg: "domain/myclass"}

	// create a Class object for the parent
	parentClass := &Class{obj: &obj{id: pid, pos: pos}}

	// Insert the parent Class object into the repository
	if err := dm.ObjRepo.Insert(parentClass); err != nil {
		t.Fatalf("Error inserting parent Class object: %v", err)
	}

	// call handleAttribute method
	dm.handleAttribute(id, pos, pid)

	// check if the repository was called with the correct Attribute object
	if len(repo.data) != 2 {
		t.Errorf("Expected repository to have 2 object, but got %v", len(repo.data))
	}

	// Check if the attribute was correctly appended to the parent Class
	if len(parentClass.attrs) != 1 {
		t.Errorf("Expected parent Class to have 1 attribute, but got %v", len(parentClass.attrs))
	}

	// call handleAttribute method with an error-prone repository
	dm.handleAttribute(id, pos, pid)

	// check if an error was pushed to the errors slice
	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
	}

	// check if the repository was not called (no object inserted)
	if len(repo.data) != 2 {
		t.Errorf("Expected repository to still have 2 object, but got %v", len(repo.data))
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

	dm.handleGenObj(id, pos)
	// check if an error was pushed to the errors slice
	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
	}

	// check if the repository was not called (no object inserted)
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to still have 1 object, but got %v", len(repo.data))
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
	dm.handleFunc(id, pos, nil, nil)

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

	dm.handleFunc(id, pos, nil, nil)
	// check if an error was pushed to the errors slice
	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
	}

	// check if the repository was not called (no object inserted)
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to still have 1 object, but got %v", len(repo.data))
	}
}

func TestDomainModel_HandleFuncWithParent(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifiers and positions
	id := &ident{name: "myFunction", pkg: "domain/myFunction"}
	pos := &pos{filename: "myFunction.go", offset: 10, line: 5, column: 15}
	pid := &ident{name: "MyClass", pkg: "domain/myclass"}

	// create a Class object for the parent
	parentClass := &Class{obj: &obj{id: pid, pos: pos}}

	// Insert the parent Class object into the repository
	if err := dm.ObjRepo.Insert(parentClass); err != nil {
		t.Fatalf("Error inserting parent Class object: %v", err)
	}

	// call handleFunc method
	dm.handleFunc(id, pos, pid, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 2 {
		t.Errorf("Expected repository to have 2 objects, but got %v", len(repo.data))
	}

	// Check if the method was correctly appended to the parent Class
	if len(parentClass.methods) != 1 {
		t.Errorf("Expected parent Class to have 1 method, but got %v", len(parentClass.methods))
	}
}

func TestDomainModel_HandleInterface(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifier and position
	id := &ident{name: "MyInterface", pkg: "domain/myinterface"}
	pos := &pos{filename: "myinterface.go", offset: 10, line: 5, column: 15}

	// call handleInterface method
	dm.handleInterface(id, pos)

	// check if the repository was called with the correct object
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to have 1 object, but got %v", len(repo.data))
	}

	var obj arch.Object
	for _, o := range repo.data {
		obj = o
		break
	}

	if _, ok := obj.(*Interface); !ok {
		t.Errorf("Expected object in repository to be an Interface, but got %T", obj)
	}

	iface := obj.(*Interface)
	if iface.obj.id.name != id.name || iface.obj.id.pkg != id.pkg {
		t.Errorf("Expected object in repository to have id %v, but got %v", id, iface.obj.id)
	}

	if iface.obj.pos.filename != pos.filename ||
		iface.obj.pos.offset != pos.offset || iface.obj.pos.line != pos.line ||
		iface.obj.pos.column != pos.column {
		t.Errorf("Expected object in repository to have pos %v, but got %v", pos, iface.obj.pos)
	}

	// call handleInterface method with an error-prone repository
	dm.handleInterface(id, pos)

	// check if an error was pushed to the errors slice
	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
	}

	// check if the repository was not called (no object inserted)
	if len(repo.data) != 1 {
		t.Errorf("Expected repository to still have 1 object, but got %v", len(repo.data))
	}
}

func TestDomainModel_HandleInterfaceMethod(t *testing.T) {
	// mock Repository
	repo := newMockRepository()

	// create Handler
	dm := &CodeHandler{Scope: "Test Model", ObjRepo: repo}

	// create identifiers and positions
	id := &ident{name: "MyMethod", pkg: "domain/mymethod"}
	pos := &pos{filename: "mymethod.go", offset: 10, line: 5, column: 15}
	pid := &ident{name: "MyInterface", pkg: "domain/myinterface"}

	// create an Interface object for the parent
	parentInterface := &Interface{obj: &obj{id: pid, pos: pos}}

	// Insert the parent Interface object into the repository
	if err := dm.ObjRepo.Insert(parentInterface); err != nil {
		t.Fatalf("Error inserting parent Interface object: %v", err)
	}

	// call handleInterfaceMethod method
	dm.handleInterfaceMethod(id, pos, pid)

	// check if the repository was called with the correct InterfaceMethod object
	if len(repo.data) != 2 {
		t.Errorf("Expected repository to have 2 objects, but got %v", len(repo.data))
	}

	// Check if the InterfaceMethod was correctly appended to the parent Interface
	if len(parentInterface.methods) != 1 {
		t.Errorf("Expected parent Interface to have 1 method, but got %v", len(parentInterface.methods))
	}

	// call handleInterfaceMethod method with an error-prone repository
	dm.handleInterfaceMethod(id, pos, pid)

	// check if an error was pushed to the errors slice
	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
	}

	// check if the repository was not called (no object inserted)
	if len(repo.data) != 2 {
		t.Errorf("Expected repository to still have 2 objects, but got %v", len(repo.data))
	}
}
func TestDomainModel_LinkHandler(t *testing.T) {
	// mock Repository
	repo := newMockRelationRepository()

	// create Handler
	dm := &CodeHandler{Scope: "test", RelRepo: repo}

	t.Run("Association Link", func(t *testing.T) {
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
			t.Errorf("Expected repository to have Association relation %v, but got %v", fromId, rel)
		}

		if obj, ok := rel.(*Association); !ok {
			t.Errorf("Expected object in repository to be an Association, but got %T", obj)
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
	})

	t.Run("Implementation Link", func(t *testing.T) {
		fromId := &ident{name: "myStruct", pkg: "/test/myStruct"}
		fromPos := &pos{filename: "myStruct.go", offset: 10, line: 5, column: 15}
		toId := &ident{name: "myInterface", pkg: "/test/myInterface"}

		link := &code.Link{
			From: &code.Node{
				Meta:   newDummyMetaWithIdent(fromId),
				Pos:    fromPos,
				Type:   code.TypeAny,
				Parent: nil,
			},
			To: &code.Node{
				Meta: newDummyMetaWithIdent(toId),
				Type: code.TypeGenInterface,
			},
			Relation: code.OneOne,
		}

		// call LinkHandler method
		dm.LinkHandler(link)

		var rel arch.Relation
		if rel = repo.Find(fromId); rel == nil {
			t.Errorf("Expected repository to have Implementation relation %v, but got %v", fromId, rel)
		}

		if obj, ok := rel.(*Implementation); !ok {
			t.Errorf("Expected object in repository to be an Implementation, but got %T", obj)
		}

		relation := rel.(*Implementation)
		if relation.From().Identifier().Name() != fromId.name || relation.Implements()[0].Identifier().Name() != toId.name {
			t.Errorf("Expected object in repository to have From %v and To %v, but got From %v and To %v",
				fromId, toId, relation.From().Identifier().Name(), relation.Implements()[0].Identifier().Name())
		}

		if relation.From().Position().Filename() != fromPos.Filename() ||
			relation.From().Position().Offset() != fromPos.Offset() ||
			relation.From().Position().Line() != fromPos.Line() ||
			relation.From().Position().Column() != fromPos.Column() {
			t.Errorf("Expected object in repository to have pos %v, but got %v", fromPos, relation.From().Position())
		}
	})

	t.Run("Dependence Link", func(t *testing.T) {
		fromId := &ident{name: "func1", pkg: "/test/func1"}
		fromPos := &pos{filename: "func1.go", offset: 10, line: 5, column: 15}
		toId := &ident{name: "func2", pkg: "/test/func2"}

		link := &code.Link{
			From: &code.Node{
				Meta:   newDummyMetaWithIdent(fromId),
				Pos:    fromPos,
				Type:   code.TypeFunc,
				Parent: nil,
			},
			To: &code.Node{
				Meta: newDummyMetaWithIdent(toId),
				Type: code.TypeFunc,
			},
			Relation: code.OneOne,
		}

		// call LinkHandler method
		dm.LinkHandler(link)

		var rel arch.Relation
		if rel = repo.Find(fromId); rel == nil {
			t.Errorf("Expected repository to have Dependence relation %v, but got %v", fromId, rel)
		}

		if obj, ok := rel.(*Dependence); !ok {
			t.Errorf("Expected object in repository to be a Dependence, but got %T", obj)
		}

		relation := rel.(*Dependence)
		if relation.From().Identifier().Name() != fromId.name || relation.DependsOn().Identifier().Name() != toId.name {
			t.Errorf("Expected object in repository to have From %v and To %v, but got From %v and To %v",
				fromId, toId, relation.From().Identifier().Name(), relation.DependsOn().Identifier().Name())
		}

		if relation.From().Position().Filename() != fromPos.Filename() ||
			relation.From().Position().Offset() != fromPos.Offset() ||
			relation.From().Position().Line() != fromPos.Line() ||
			relation.From().Position().Column() != fromPos.Column() {
			t.Errorf("Expected object in repository to have pos %v, but got %v", fromPos, relation.From().Position())
		}
	})

	t.Run("Composition Link", func(t *testing.T) {
		fromId := &ident{name: "struct1", pkg: "/test/struct1"}
		fromPos := &pos{filename: "struct1.go", offset: 10, line: 5, column: 15}
		toId := &ident{name: "func1", pkg: "/test/func1"}

		link := &code.Link{
			From: &code.Node{
				Meta:   newDummyMetaWithIdent(fromId),
				Pos:    fromPos,
				Type:   code.TypeGenStruct,
				Parent: nil,
			},
			To: &code.Node{
				Meta: newDummyMetaWithIdent(toId),
				Type: code.TypeGenStructField,
			},
			Relation: code.OneOne,
		}

		// call LinkHandler method
		dm.LinkHandler(link)

		var rel arch.Relation
		if rel = repo.Find(fromId); rel == nil {
			t.Errorf("Expected repository to have Composition relation %v, but got %v", fromId, rel)
		}

		if obj, ok := rel.(*Composition); !ok {
			t.Errorf("Expected object in repository to be a Composition, but got %T", obj)
		}

		relation := rel.(*Composition)
		if relation.From().Identifier().Name() != fromId.name || relation.Child().Identifier().Name() != toId.name {
			t.Errorf("Expected object in repository to have From %v and To %v, but got From %v and To %v",
				fromId, toId, relation.From().Identifier().Name(), relation.Child().Identifier().Name())
		}

		if relation.From().Position().Filename() != fromPos.Filename() ||
			relation.From().Position().Offset() != fromPos.Offset() ||
			relation.From().Position().Line() != fromPos.Line() ||
			relation.From().Position().Column() != fromPos.Column() {
			t.Errorf("Expected object in repository to have pos %v, but got %v", fromPos, relation.From().Position())
		}
	})

	t.Run("Embedding Link", func(t *testing.T) {
		repo.Clear()

		fromId := &ident{name: "struct1", pkg: "/test/struct1"}
		fromPos := &pos{filename: "struct1.go", offset: 10, line: 5, column: 15}
		toId := &ident{name: "field1", pkg: "/test/field1"}

		link := &code.Link{
			From: &code.Node{
				Meta:   newDummyMetaWithIdent(fromId),
				Pos:    fromPos,
				Type:   code.TypeGenStruct,
				Parent: nil,
			},
			To: &code.Node{
				Meta: newDummyMetaWithIdent(toId),
				Type: code.TypeGenStructEmbeddedField,
			},
			Relation: code.OneOne,
		}

		// call LinkHandler method
		dm.LinkHandler(link)

		var rel arch.Relation
		if rel = repo.Find(fromId); rel == nil {
			t.Errorf("Expected repository to have Embedding relation %v, but got %v", fromId, rel)
		}

		if obj, ok := rel.(*Embedding); !ok {
			t.Errorf("Expected object in repository to be an Embedding, but got %T", obj)
		}

		relation := rel.(*Embedding)
		if relation.From().Identifier().Name() != fromId.name || relation.Embedded().Identifier().Name() != toId.name {
			t.Errorf("Expected object in repository to have From %v and To %v, but got From %v and To %v",
				fromId, toId, relation.From().Identifier().Name(), relation.Embedded().Identifier().Name())
		}

		if relation.From().Position().Filename() != fromPos.Filename() ||
			relation.From().Position().Offset() != fromPos.Offset() ||
			relation.From().Position().Line() != fromPos.Line() ||
			relation.From().Position().Column() != fromPos.Column() {
			t.Errorf("Expected object in repository to have pos %v, but got %v", fromPos, relation.From().Position())
		}
	})
}

func TestDomainModel_LinkHandler_OutOfScope(t *testing.T) {
	// mock Repository
	repo := newMockRelationRepository()

	// create Handler with a specific scope
	dm := &CodeHandler{Scope: "test", RelRepo: repo}

	// create link with source and target outside of the specified scope
	fromId := &ident{name: "field", pkg: "/other/pkg"}
	fromPos := &pos{filename: "other_pkg.go", offset: 10, line: 5, column: 15}
	toId := &ident{name: "func", pkg: "/yet/another/pkg"}

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

	// Check that no relation was inserted into the repository
	if len(repo.relations) != 0 {
		t.Errorf("Expected no relation to be inserted, but got %v", len(repo.relations))
	}
}

func TestDomainModel_LinkHandler_ErrorInsert(t *testing.T) {
	// mock Repository with an error-prone Insert method
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
	dm.LinkHandler(link)

	// Check that an error was pushed to the errors slice
	if len(dm.errors) != 1 {
		t.Errorf("Expected 1 error, but got %v", len(dm.errors))
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
			name: "Test handle struct field for TypeGenStructField with parent nil",
			node: &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeGenStructField,
				Parent: nil},
			want:     nil,
			wantType: nil,
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
			name: "Test handle interface method for TypeGenInterfaceMethod with parent nil",
			node: &code.Node{Meta: newDummyMetaWithIdent(id1), Pos: pos1, Type: code.TypeGenInterfaceMethod,
				Parent: nil},
			want:     nil,
			wantType: nil,
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

			if got != nil && tt.want != nil && got.Identifier().ID() != tt.want.Identifier().ID() {
				t.Errorf("NodeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeHandler_pushError(t *testing.T) {
	ch := &CodeHandler{}
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	ch.pushError(err1)
	ch.pushError(err2)

	// 检查是否添加了正确数量的错误
	if len(ch.errors) != 2 {
		t.Errorf("Expected 2 errors, but got %d", len(ch.errors))
	}

	// 检查是否添加了正确的错误
	if ch.errors[0] != err1 {
		t.Errorf("Expected error 1, but got: %v", ch.errors[0])
	}

	if ch.errors[1] != err2 {
		t.Errorf("Expected error 2, but got: %v", ch.errors[1])
	}
}
