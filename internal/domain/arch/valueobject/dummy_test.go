package valueobject

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
)

type DummyMeta struct {
	pkg        string
	name       string
	parentName string
}

func (m *DummyMeta) Pkg() string     { return m.pkg }
func (m *DummyMeta) Name() string    { return m.name }
func (m *DummyMeta) Parent() string  { return m.parentName }
func (m *DummyMeta) HasParent() bool { return m.parentName != "" }

func newDummyMetaWithIdent(id *ident) *DummyMeta {
	return &DummyMeta{
		pkg:        id.pkg,
		name:       id.name,
		parentName: "", // Set parentName based on your requirements
	}
}

type MockPosition struct {
	filename string
	offset   int
	line     int
	column   int
}

func (p *MockPosition) Filename() string { return p.filename }
func (p *MockPosition) Offset() int      { return p.offset }
func (p *MockPosition) Line() int        { return p.line }
func (p *MockPosition) Column() int      { return p.column }

type DummyMetaInfo struct {
	pkg        string
	name       string
	parentName string
}

func (m *DummyMetaInfo) Pkg() string     { return m.pkg }
func (m *DummyMetaInfo) Name() string    { return m.name }
func (m *DummyMetaInfo) Parent() string  { return m.parentName }
func (m *DummyMetaInfo) HasParent() bool { return m.parentName != "" }

func newMockRepository() *mockObjRepository {
	return &mockObjRepository{data: map[arch.ObjIdentifier]arch.Object{}}
}

type mockObjRepository struct {
	data map[arch.ObjIdentifier]arch.Object
}

func (r *mockObjRepository) Clear() {
	r.data = map[arch.ObjIdentifier]arch.Object{}
}

func (r *mockObjRepository) Find(id arch.ObjIdentifier) arch.Object {
	for key, value := range r.data {
		if key.ID() == id.ID() {
			return value
		}
	}
	return nil
}

func (r *mockObjRepository) Insert(obj arch.Object) error {
	r.data[obj.Identifier()] = obj
	return nil
}

func (r *mockObjRepository) Walk(cb func(obj arch.Object) error) {
	for _, obj := range r.data {
		if err := cb(obj); err != nil {
			return
		}
	}
}

func (r *mockObjRepository) GetObjects(ids []arch.ObjIdentifier) ([]arch.Object, error) {
	var objs []arch.Object
	for _, id := range ids {
		objs = append(objs, r.data[id])
	}
	return objs, nil
}

func (r *mockObjRepository) All() []arch.ObjIdentifier {
	var ids []arch.ObjIdentifier
	for i := range r.data {
		ids = append(ids, i)
	}
	return ids
}

func newMockRelationRepository() *mockRelationRepository {
	return &mockRelationRepository{
		relations: make([]arch.Relation, 0),
	}
}

type mockRelationRepository struct {
	relations []arch.Relation
}

func (r *mockRelationRepository) Find(id arch.Identifier) arch.Relation {
	for _, rel := range r.relations {
		if rel.From().Identifier().ID() == id.ID() {
			return rel
		}
	}
	return nil
}

func (r *mockRelationRepository) Insert(rel arch.Relation) error {
	r.relations = append(r.relations, rel)
	return nil
}

func (r *mockRelationRepository) Walk(walker func(rel arch.Relation) error) {
	for _, rel := range r.relations {
		err := walker(rel)
		if err != nil {
			fmt.Printf("Error while walking relations: %v\n", err)
			break
		}
	}
}

type MockObject struct {
	id                     arch.ObjIdentifier
	position               arch.Position
	name                   string
	dir                    string
	NameSeparatorLengthVal int
}

// ID returns the identifier of the object
func (mo MockObject) ID() string {
	return mo.id.ID()
}

func (mo MockObject) Identifier() arch.ObjIdentifier {
	return mo.id
}

func (mo MockObject) Position() arch.Position {
	return mo.position
}

// Name returns the name of the object
func (mo MockObject) Name() string {
	return mo.name
}

// Dir returns the directory of the object
func (mo MockObject) Dir() string {
	return mo.dir
}

func (mo MockObject) NameSeparatorLength() int {
	return mo.NameSeparatorLengthVal
}
