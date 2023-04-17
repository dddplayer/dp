package factory

import (
	"github.com/dddplayer/core/codeanalysis/entity"
)

func NewPkg(path string) (*entity.Pkg, error) {
	p := &entity.Pkg{Path: path}
	if err := p.Load(); err != nil {
		return nil, err
	}
	return p, nil
}
