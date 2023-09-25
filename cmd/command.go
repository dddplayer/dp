package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/interfaces/cmd"
	"os"
)

func New() error {
	topLevel := flag.NewFlagSet("dp", flag.ExitOnError)
	topLevel.Usage = func() {
		fmt.Println("Usage:\n  dp [command]")
		fmt.Println("\nCommands:")
		fmt.Println("  strategic:  generate domain strategic diagram")
		fmt.Println("  tactic:     generate domain tactic diagram")
		fmt.Println("  normal:     generate normal arch diagram")

		fmt.Println("\nExample:")
		fmt.Println("  dp normal -m ~/github/dddplayer/dp -p github.com/dddplayer/dp/internal/domain")
	}

	err := topLevel.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if topLevel.Parsed() {
		if len(topLevel.Args()) == 0 {
			topLevel.Usage()
			return errors.New("please specify a sub-command")
		}

		// 获取子命令及参数
		subCommand := topLevel.Args()[0]

		switch subCommand {
		case "normal":
			normalCmd, err := cmd.NewNormalCmd(topLevel)
			if err != nil {
				return err
			}
			if err := normalCmd.Run(); err != nil {
				return err
			}

		case "strategic":
			strategicCmd, err := cmd.NewStrategicCmd(topLevel)
			if err != nil {
				return err
			}
			if err := strategicCmd.Run(); err != nil {
				return err
			}

		case "tactic":
			tacticCmd, err := cmd.NewTacticCmd(topLevel)
			if err != nil {
				return err
			}
			if err := tacticCmd.Run(); err != nil {
				return err
			}

		default:
			topLevel.Usage()
			return errors.New("invalid sub-command")
		}
	}

	return nil
}
