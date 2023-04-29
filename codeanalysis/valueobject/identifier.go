package valueobject

import "fmt"

const (
	NodeJoiner = "."
)

type Identifier interface {
	Path() string
	Name() string
	String() string
}

func NewIdentifier(path, name string) Identifier {
	return &identifier{
		path: path,
		name: name,
	}
}

type identifier struct {
	path string
	name string
}

func (i *identifier) String() string {
	return fmt.Sprintf("%s%s%s", i.path, NodeJoiner, i.name)
}

func (i *identifier) Path() string {
	return i.path
}

func (i *identifier) Name() string {
	return i.name
}
