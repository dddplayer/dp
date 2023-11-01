package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/application"
	"github.com/dddplayer/dp/internal/infrastructure/persistence"
	"strings"
)

type normalCmd struct {
	parent     *flag.FlagSet
	cmd        *flag.FlagSet
	mainFlag   *string
	pkgFlag    *string
	comFlag    *bool
	detailFlag *bool
	mfFlag     *bool
}

func NewNormalCmd(parent *flag.FlagSet) (*normalCmd, error) {
	nCmd := &normalCmd{
		parent: parent,
	}

	nCmd.cmd = flag.NewFlagSet("normal", flag.ExitOnError)
	nCmd.mainFlag = nCmd.cmd.String("m", "", fmt.Sprintf(
		"[required] main package path \n(e.g. %s)", "github.com/dddplayer/dp"))
	nCmd.pkgFlag = nCmd.cmd.String("p", "", fmt.Sprintf(
		"[required] target package path \n(e.g. %s)", "github.com/dddplayer/dp/internal/domain"))
	nCmd.comFlag = nCmd.cmd.Bool("c", false, "show struct composition relation")
	nCmd.detailFlag = nCmd.cmd.Bool("d", false, "show all relations")
	nCmd.mfFlag = nCmd.cmd.Bool("mf", false, "show message flow relations")

	err := nCmd.cmd.Parse(parent.Args()[1:])
	if err != nil {
		return nil, err
	}

	return nCmd, nil
}

func (nc *normalCmd) Usage() {
	nc.cmd.Usage()
}

func (nc *normalCmd) Run() error {
	if *nc.mainFlag == "" {
		nc.cmd.Usage()
		return errors.New("please specify the main package")
	}

	if *nc.pkgFlag == "" {
		nc.cmd.Usage()
		return errors.New("please specify a target package full name")
	}

	if *nc.comFlag {
		return normalCompositionGraph(*nc.mainFlag, *nc.pkgFlag)
	}

	if *nc.mfFlag {
		return normalMessageFlowGraph(*nc.mainFlag, *nc.pkgFlag)
	}

	if *nc.detailFlag {
		return normalDetailGraph(*nc.mainFlag, *nc.pkgFlag)
	}

	return normalGraph(*nc.mainFlag, *nc.pkgFlag)
}

func normalCompositionGraph(mainPkg, domain string) error {
	dot, err := application.CompositionGeneralGraph(mainPkg, domain,
		persistence.NewRadixTree(),
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	if err = open(dot); err != nil {
		return err
	}

	if err = writeToDisk(dot, filename(domain, "composition"), mainPkg); err != nil {
		return err
	}

	return nil
}

func normalDetailGraph(mainPkg, domain string) error {
	dot, err := application.DetailGeneralGraph(mainPkg, domain,
		persistence.NewRadixTree(),
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	if err = open(dot); err != nil {
		return err
	}

	if err = writeToDisk(dot, filename(domain, "detail"), mainPkg); err != nil {
		return err
	}

	return nil
}

func normalMessageFlowGraph(mainPkg, domain string) error {
	dot, err := application.MessageFlowGraph(mainPkg, domain,
		persistence.NewRadixTree(),
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	if err = open(dot); err != nil {
		return err
	}

	if err = writeToDisk(dot, filename(domain, "messageflow"), mainPkg); err != nil {
		return err
	}

	return nil
}

func normalGraph(mainPkg, domain string) error {
	dot, err := application.GeneralGraph(mainPkg, domain,
		persistence.NewRadixTree(),
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	if err = open(dot); err != nil {
		return err
	}

	if err = writeToDisk(dot, filename(domain, ""), mainPkg); err != nil {
		return err
	}

	return nil
}

func filename(main, sub string) string {
	mainStr := strings.ReplaceAll(main, "/", ".")
	if sub == "" {
		return mainStr
	}
	return fmt.Sprintf("%s.%s", mainStr, sub)
}
