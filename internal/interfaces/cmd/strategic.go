package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/application"
	"github.com/dddplayer/dp/internal/infrastructure/persistence"
	"path"
	"strings"
)

type strategicCmd struct {
	parent   *flag.FlagSet
	cmd      *flag.FlagSet
	mainFlag *string
	pkgFlag  *string
}

func NewStrategicCmd(parent *flag.FlagSet) (*strategicCmd, error) {
	sCmd := &strategicCmd{
		parent: parent,
	}

	sCmd.cmd = flag.NewFlagSet("strategic", flag.ExitOnError)
	sCmd.mainFlag = sCmd.cmd.String("m", "", fmt.Sprintf(
		"[required] main package path \n(e.g. %s)", "github.com/dddplayer/dp"))
	sCmd.pkgFlag = sCmd.cmd.String("p", "", fmt.Sprintf(
		"[required] target package \n(e.g. %s)", "github.com/dddplayer/dp/internal/domain"))

	err := sCmd.cmd.Parse(parent.Args()[1:])
	if err != nil {
		return nil, err
	}

	return sCmd, nil
}

func (sc *strategicCmd) Usage() {
	sc.cmd.Usage()
}

func (sc *strategicCmd) Run() error {
	if *sc.mainFlag == "" {
		sc.cmd.Usage()
		return errors.New("please specify the main package")
	}

	if *sc.pkgFlag == "" {
		sc.cmd.Usage()
		return errors.New("please specify a target package full name")
	}

	return strategicGraph(*sc.mainFlag, *sc.pkgFlag)
}

func strategicGraph(mainPkg, domain string) error {
	dot, err := application.StrategicGraph(mainPkg, domain,
		persistence.NewRadixTree(),
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	open(dot)
	if err = writeToDisk(dot, strings.ReplaceAll(path.Join(domain, "detail"), "/", "."), mainPkg); err != nil {
		return err
	}

	return nil
}
