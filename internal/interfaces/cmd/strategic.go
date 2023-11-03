package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/application"
	"github.com/dddplayer/dp/internal/infrastructure/persistence"
)

type strategicCmd struct {
	parent       *flag.FlagSet
	cmd          *flag.FlagSet
	mainFlag     *string
	pkgFlag      *string
	fastModeFlag *bool
	deepModeFlag *bool
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
	sCmd.fastModeFlag = sCmd.cmd.Bool("fast", true, "analysis code in fast mode to save time")
	sCmd.deepModeFlag = sCmd.cmd.Bool("deep", false, "analysis code in fast mode to get more accurate information")

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

	if *sc.deepModeFlag {
		return strategicGraph(*sc.mainFlag, *sc.pkgFlag, true)
	}

	return strategicGraph(*sc.mainFlag, *sc.pkgFlag, false)
}

func strategicGraph(mainPkg, domain string, deep bool) error {
	dot, err := application.StrategicGraph(mainPkg, domain, deep,
		persistence.NewRadixTree(),
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	if err = open(dot); err != nil {
		return err
	}
	if err = writeToDisk(dot, filename(domain, "strategic"), mainPkg); err != nil {
		return err
	}

	return nil
}
