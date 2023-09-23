package cmd

import (
	"flag"
	"fmt"
	"github.com/dddplayer/dp/internal/application"
	"github.com/dddplayer/dp/internal/infrastructure/persistence"
	"github.com/dddplayer/dp/pkg/datastructure/radix"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

var dotCmd = flag.NewFlagSet("dot", flag.ExitOnError)
var dotGeneral, dotGeneralDetail, dotStrategic, dotTactic, dotTacticDetail *bool

func DotUsage() {
	dotCmd.Usage = func() {
		fmt.Println("  dp [-m Main package] [-d Domain module] dot [-g General | -gd General Detail -s Strategic | -t Tactic | -td Tactic Detail]")
	}
	dotCmd.Usage()
}

func DotCmd() *flag.FlagSet {
	dotGeneral = dotCmd.Bool("g", false, "general dot directed")
	dotGeneralDetail = dotCmd.Bool("gd", false, "general detail dot directed")
	dotStrategic = dotCmd.Bool("s", false, "strategic dot directed")
	dotTactic = dotCmd.Bool("t", false, "tactic dot directed")
	dotTacticDetail = dotCmd.Bool("td", false, "tactic detail dot directed")
	return dotCmd
}

func DotRun(mainPkg, domain string) error {
	if *dotGeneral {
		dot, err := application.GeneralGraph(mainPkg, domain,
			&persistence.RadixTree{Tree: radix.NewTree()},
			&persistence.Relations{}, false,
		)
		if err != nil {
			return err
		}

		open(dot)
		writeToDisk(dot, strings.ReplaceAll(domain, "/", "."))
	}

	if *dotGeneralDetail {
		dot, err := application.GeneralGraph(mainPkg, domain,
			&persistence.RadixTree{Tree: radix.NewTree()},
			&persistence.Relations{}, true,
		)
		if err != nil {
			return err
		}

		open(dot)
		writeToDisk(dot, strings.ReplaceAll(domain, "/", "."))
	}

	if *dotStrategic {
		dot, err := application.StrategicGraph(mainPkg, domain,
			&persistence.RadixTree{Tree: radix.NewTree()},
			&persistence.Relations{},
		)
		if err != nil {
			return err
		}

		open(dot)
	}
	if *dotTactic {
		dot := application.TacticGraph(mainPkg, domain,
			&persistence.RadixTree{Tree: radix.NewTree()},
			&persistence.Relations{},
			false,
		)

		open(dot)
		writeToDisk(dot, strings.ReplaceAll(domain, "/", "."))
	}
	if *dotTacticDetail {
		dot := application.TacticGraph(mainPkg, domain,
			&persistence.RadixTree{Tree: radix.NewTree()},
			&persistence.Relations{},
			true,
		)

		open(dot)
		writeToDisk(dot, strings.ReplaceAll(path.Join(domain, "detail"), "/", "."))
	}

	return nil
}

func writeToDisk(raw string, filename string) {
	// Open file for writing
	file, err := os.Create(fmt.Sprintf("%s.dot", filename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write byte slice to file
	_, err = file.Write([]byte(raw))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("String written to file successfully.")
}

func open(raw string) {
	encoded := encodeURIComponent(raw)
	err := openBrowser(fmt.Sprintf("https://dddplayer.com/#%s", encoded))
	if err != nil {
		return
	}
}

func encodeURIComponent(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

func openBrowser(url string) error {
	cmd := exec.Command("open", url)
	return cmd.Start()
}
