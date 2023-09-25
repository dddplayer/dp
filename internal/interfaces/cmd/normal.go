package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/application"
	"github.com/dddplayer/dp/internal/infrastructure/persistence"
	"github.com/dddplayer/dp/pkg/datastructure/radix"
	"strings"
)

type normalCmd struct {
	parent     *flag.FlagSet
	cmd        *flag.FlagSet
	mainFlag   *string
	pkgFlag    *string
	comFlag    *bool
	detailFlag *bool
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

	if *nc.detailFlag {
		return normalDetailGraph(*nc.mainFlag, *nc.pkgFlag)
	}

	return normalGraph(*nc.mainFlag, *nc.pkgFlag)
}

func normalCompositionGraph(mainPkg, domain string) error {
	dot, err := application.CompositionGeneralGraph(mainPkg, domain,
		&persistence.RadixTree{Tree: radix.NewTree()},
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	open(dot)
	writeToDisk(dot, strings.ReplaceAll(domain, "/", "."))

	return nil
}

func normalDetailGraph(mainPkg, domain string) error {
	dot, err := application.DetailGeneralGraph(mainPkg, domain,
		&persistence.RadixTree{Tree: radix.NewTree()},
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	open(dot)
	writeToDisk(dot, strings.ReplaceAll(domain, "/", "."))

	return nil
}

func normalGraph(mainPkg, domain string) error {
	dot, err := application.GeneralGraph(mainPkg, domain,
		&persistence.RadixTree{Tree: radix.NewTree()},
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	open(dot)
	writeToDisk(dot, strings.ReplaceAll(domain, "/", "."))

	return nil
}
