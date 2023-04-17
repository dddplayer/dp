package service

import (
	"github.com/dddplayer/core/codeanalysis/service"
	"github.com/dddplayer/core/factory"
)

func Dot(domain string) string {
	dm, err := factory.NewDomainModel()
	if err != nil {
		return err.Error()
	}
	if err := service.Visit(domain, dm.NameHandler); err != nil {
		return err.Error()
	}

	d, err := dm.Output()
	if err != nil {
		return ""
	}

	return d
}
