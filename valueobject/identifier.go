package valueobject

import (
	"fmt"
	"github.com/dddplayer/core/codeanalysis/valueobject"
	"github.com/dddplayer/core/dot/entity"
	"path"
)

const (
	DomainObjJoiner  = valueobject.NodeJoiner
	DomainPortJoiner = entity.DotPortJoiner
)

type Identifier struct {
	Name string
	Path string
}

func (i *Identifier) String() string {
	return fmt.Sprintf("%s%s%s", i.Path, DomainObjJoiner, i.Name)
}

func (i *Identifier) DomainName() string {
	return path.Base(path.Dir(i.Path))
}

func (i *Identifier) Base() string {
	return path.Base(i.String())
}

func NewIdentifier(id valueobject.Identifier) Identifier {
	return Identifier{
		Name: id.Name(),
		Path: id.Path(),
	}
}
