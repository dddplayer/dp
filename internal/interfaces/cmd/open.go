package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type openCmd struct {
	parent          *flag.FlagSet
	cmd             *flag.FlagSet
	archDiagramPath *string
}

func NewOpenCmd(parent *flag.FlagSet) (*openCmd, error) {
	nCmd := &openCmd{
		parent: parent,
	}

	nCmd.cmd = flag.NewFlagSet("normal", flag.ExitOnError)
	nCmd.archDiagramPath = nCmd.cmd.String("p", "", fmt.Sprintf(
		"[required] target arch diagram path \n(e.g. %s)", "dddplayer/arch.dot"))

	err := nCmd.cmd.Parse(parent.Args()[1:])
	if err != nil {
		return nil, err
	}

	return nCmd, nil
}

func (oc *openCmd) Usage() {
	oc.cmd.Usage()
}

func (oc *openCmd) Run() error {
	if *oc.archDiagramPath == "" {
		oc.cmd.Usage()
		return errors.New("please specify a target arch diagram path")
	}

	_, err := os.Stat(*oc.archDiagramPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", *oc.archDiagramPath)
	}

	dotStr, err := os.ReadFile(*oc.archDiagramPath)
	if err != nil {
		return err
	}

	return open(string(dotStr))
}
