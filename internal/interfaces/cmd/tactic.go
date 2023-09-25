package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/application"
	"github.com/dddplayer/dp/internal/infrastructure/persistence"
	"github.com/dddplayer/dp/pkg/datastructure/radix"
	"path"
	"strings"
)

type tacticCmd struct {
	parent     *flag.FlagSet
	cmd        *flag.FlagSet
	mainFlag   *string
	pkgFlag    *string
	detailFlag *bool
}

func NewTacticCmd(parent *flag.FlagSet) (*tacticCmd, error) {
	tCmd := &tacticCmd{
		parent: parent,
	}

	tCmd.cmd = flag.NewFlagSet("strategic", flag.ExitOnError)
	tCmd.mainFlag = tCmd.cmd.String("m", "", fmt.Sprintf(
		"[required] main package path \n(e.g. %s)", "github.com/dddplayer/dp"))
	tCmd.pkgFlag = tCmd.cmd.String("p", "", fmt.Sprintf(
		"[required] target package \n(e.g. %s)", "github.com/dddplayer/dp/internal/domain"))
	tCmd.detailFlag = tCmd.cmd.Bool("d", false, "show all relations")

	err := tCmd.cmd.Parse(parent.Args()[1:])
	if err != nil {
		return nil, err
	}

	return tCmd, nil
}

func (sc *tacticCmd) Usage() {
	sc.cmd.Usage()
}

func (sc *tacticCmd) Run() error {
	if *sc.mainFlag == "" {
		sc.cmd.Usage()
		return errors.New("please specify the main package")
	}

	if *sc.pkgFlag == "" {
		sc.cmd.Usage()
		return errors.New("please specify a target package full name")
	}

	if *sc.detailFlag {
		return detailTacticGraph(*sc.mainFlag, *sc.pkgFlag)
	}

	return tacticGraph(*sc.mainFlag, *sc.pkgFlag)
}

func tacticGraph(mainPkg, domain string) error {
	dot, err := application.TacticGraph(mainPkg, domain,
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

func detailTacticGraph(mainPkg, domain string) error {
	dot, err := application.DetailTacticGraph(mainPkg, domain,
		&persistence.RadixTree{Tree: radix.NewTree()},
		&persistence.Relations{},
	)
	if err != nil {
		return err
	}

	open(dot)
	writeToDisk(dot, strings.ReplaceAll(path.Join(domain, "detail"), "/", "."))
	return nil
}
