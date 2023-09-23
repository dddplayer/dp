package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/code"
)

func NewMeta(pkg, name string) code.MetaInfo {
	return &meta{
		pkg:        pkg,
		name:       name,
		parentName: "",
	}
}

func NewMetaWithParent(pkg, name, parentName string) code.MetaInfo {
	return &meta{
		pkg:        pkg,
		name:       name,
		parentName: parentName,
	}
}

type meta struct {
	pkg        string
	name       string
	parentName string
}

func (m *meta) Pkg() string {
	return m.pkg
}

func (m *meta) Parent() string {
	return m.parentName
}

func (m *meta) Name() string {
	return m.name
}

func (m *meta) HasParent() bool {
	return m.parentName != ""
}
