package persistence

import "github.com/dddplayer/dp/internal/domain/arch"

type MockObject struct {
	id       arch.ObjIdentifier
	position arch.Position
}

func (mo MockObject) Identifier() arch.ObjIdentifier {
	return mo.id
}

func (mo MockObject) Position() arch.Position {
	return mo.position
}

type MockRelation struct {
	relationType arch.RelationType
	fromObject   arch.Object
}

func (mr *MockRelation) Type() arch.RelationType {
	return mr.relationType
}

func (mr *MockRelation) From() arch.Object {
	return mr.fromObject
}

type MockIdentifier struct {
	IDVal                  string
	NameVal                string
	NameSeparatorLengthVal int
	DirVal                 string
}

func (mi *MockIdentifier) ID() string {
	return mi.IDVal
}

func (mi *MockIdentifier) Name() string {
	return mi.NameVal
}

func (mi *MockIdentifier) NameSeparatorLength() int {
	return mi.NameSeparatorLengthVal
}

func (mi *MockIdentifier) Dir() string {
	return mi.DirVal
}
