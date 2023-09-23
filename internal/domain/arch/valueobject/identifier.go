package valueobject

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/code"
	"path"
	"strings"
)

const (
	DotJoiner = "."
)

type ident struct {
	name string
	pkg  string
}

func (i *ident) Name() string             { return i.name }
func (i *ident) NameSeparatorLength() int { return len(DotJoiner) }
func (i *ident) Dir() string              { return i.pkg }
func (i *ident) ID() string               { return path.Join(i.pkg, i.name) }
func (i *ident) fixTmpName()              { i.name = strings.Split(i.name, "$")[0] }

func newIdentifier(meta code.MetaInfo) *ident {
	name := meta.Name()
	if meta.HasParent() {
		name = fmt.Sprintf("%s%s%s", meta.Parent(), DotJoiner, name)
	}
	return &ident{
		name: name,
		pkg:  meta.Pkg(),
	}
}
