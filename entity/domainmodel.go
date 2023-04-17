package entity

import (
	"bytes"
	"fmt"
	"github.com/dddplayer/core/dot/entity"
	"github.com/dddplayer/core/dot/service"
)

type DomainModel struct {
	Name string
}

func (dm *DomainModel) NameHandler(name string) {
	dm.Name = name
}

func (dm *DomainModel) Output() (string, error) {
	g := &dotGraph{
		name:  dm.Name,
		nodes: []entity.DotNode{&dotNode{name: fmt.Sprintf("node_%s", dm.Name)}},
	}
	var buf bytes.Buffer
	if err := service.WriteDot(g, &buf); err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
