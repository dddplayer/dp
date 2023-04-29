package service

import (
	"github.com/dddplayer/core/codeanalysis/entity"
	"github.com/dddplayer/core/codeanalysis/factory"
)

func Visit(mainPkgPath, domain string, nodeCB entity.NodeCB, linkCB entity.LinkCB) error {
	p, err := factory.NewPkg(mainPkgPath, domain)
	if err != nil {
		return err
	}

	p.VisitFile(nodeCB, linkCB)
	p.InterfaceImplements(linkCB)
	if err := p.CallGraph(linkCB); err != nil {
		return err
	}

	return nil
}
