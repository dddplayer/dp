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
	mainPkg := topLevel.String("m", "", "main package")
	domainModule := topLevel.String("d", "", "domain module")
	topLevel.Usage = func() {
		fmt.Println("Usage: dp [options] [subcommand]")
		fmt.Println("  dp [-m Main package] [-d Domain module] [subcommand]")
		cmd.DotUsage()

		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  dp -m ./example -d github.com/dddplayer/markdown dot -s")
	}

	err := topLevel.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if topLevel.Parsed() {
		if *mainPkg == "" {
			topLevel.Usage()
			return errors.New("please specify the main package path")
		}

		if *domainModule == "" {
			topLevel.Usage()
			return errors.New("please specify a domain module name")
		}

		if len(topLevel.Args()) == 0 {
			topLevel.Usage()
			return errors.New("please specify a sub-command")
		}

		// 获取子命令及参数
		subCommand := topLevel.Args()[0]

		switch subCommand {
		case "dot":
			dot := cmd.DotCmd()
			err := dot.Parse(topLevel.Args()[1:])
			if err != nil {
				dot.Usage()
				return err
			}

			// 处理子命令1及参数
			if dot.Parsed() {
				err := cmd.DotRun(*mainPkg, *domainModule)
				if err != nil {
					dot.Usage()
					return err
				}
			}

		default:
			topLevel.Usage()
			return errors.New("invalid sub-command")
		}
	}

	return nil
}
