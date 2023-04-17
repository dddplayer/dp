package main

import (
	"flag"
	"fmt"
	"github.com/dddplayer/core/service"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		os.Exit(-1)
	}

	dot := service.Dot(args[0])
	fmt.Println(dot)
}
