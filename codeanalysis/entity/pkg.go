package entity

import (
	"fmt"
	"go/build"
	"golang.org/x/tools/go/packages"
)

type Pkg struct {
	Path    string
	Initial []*packages.Package
}

func (p *Pkg) Load() error {
	cfg := &packages.Config{
		Mode:       packages.LoadAllSyntax,
		Tests:      false,
		Dir:        "",
		BuildFlags: build.Default.BuildTags,
	}

	initial, err := packages.Load(cfg, p.Path)
	if err != nil {
		return err
	}
	if packages.PrintErrors(initial) > 0 {
		return fmt.Errorf("packages contain errors")
	}
	if len(initial) == 0 {
		return fmt.Errorf("package empty error")
	}

	p.Initial = initial
	return nil
}
