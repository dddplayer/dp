package factory

import (
	"github.com/dddplayer/core/codeanalysis/entity"
)

func NewPkg(path, domain string) (*entity.Pkg, error) {
	p := &entity.Pkg{Path: path, DomainPkgPath: domain}
	if err := p.Load(); err != nil {
		return nil, err
	}
	return p, nil
}
