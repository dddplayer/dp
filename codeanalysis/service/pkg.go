package service

import (
	"github.com/dddplayer/core/codeanalysis/factory"
)

func Visit(path string, name func(name string)) error {
	p, err := factory.NewPkg(path)
	if err != nil {
		return err
	}
	name(p.Initial[0].PkgPath)
	return nil
}
