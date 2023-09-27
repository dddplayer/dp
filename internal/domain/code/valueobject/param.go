package valueobject

import "github.com/dddplayer/dp/internal/domain/code"

type Params []code.Param

func (p Params) Contains(name string) bool {
	for _, param := range p {
		if param.Name() == name {
			return true
		}
	}
	return false
}

type param struct {
	name string
}

func (p *param) Name() string {
	return p.name
}

func NewParam(name string) code.Param {
	return &param{name: name}
}
