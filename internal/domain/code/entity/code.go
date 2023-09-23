package entity

import "github.com/dddplayer/dp/internal/domain/code"

type Code struct {
	lan code.Language
}

func NewCode(mainPkgPath, domain string) (*Code, error) {
	g, err := newGo(mainPkgPath, domain)
	if err != nil {
		return nil, err
	}
	return &Code{lan: g}, nil
}

func newGo(path, domain string) (code.Language, error) {
	p := &Go{Path: path, DomainPkgPath: domain}
	if err := p.Load(); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Code) VisitFast(handler code.Handler) error {
	c.lan.VisitFile(handler.NodeHandler, handler.LinkHandler)
	c.lan.InterfaceImplements(handler.LinkHandler)
	if err := c.lan.CallGraph(handler.LinkHandler, code.CallGraphFastMode); err != nil {
		return err
	}

	return nil
}

func (c *Code) VisitDeep(handler code.Handler) error {
	c.lan.VisitFile(handler.NodeHandler, handler.LinkHandler)
	c.lan.InterfaceImplements(handler.LinkHandler)
	if err := c.lan.CallGraph(handler.LinkHandler, code.CallGraphDeepMode); err != nil {
		return err
	}

	return nil
}
