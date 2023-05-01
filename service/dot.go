package service

import (
	"bytes"
	"github.com/dddplayer/core/codeanalysis/service"
	dot "github.com/dddplayer/core/dot/service"
	"github.com/dddplayer/core/entity"
	"github.com/dddplayer/core/factory"
)

func Dot(mainPkgPath string, domain string, repo entity.Repository) string {
	dm, err := factory.NewDomainModel(domain, repo)
	if err != nil {
		return err.Error()
	}
	if err := service.Visit(mainPkgPath, domain, dm.NodeHandler, dm.LinkHandler); err != nil {
		return err.Error()
	}

	db := &factory.DotBuilder{
		Repo:   dm.Repo,
		Domain: entity.NewDomain(dm.Name),
	}
	g, err := db.Build()
	if err != nil {
		return err.Error()
	}

	var buf bytes.Buffer
	if err := dot.WriteDot(g, &buf); err != nil {
		return err.Error()
	}

	return string(buf.Bytes())
}
