package main

import (
	"flag"
	"fmt"
	"github.com/dddplayer/core/datastructure/radix"
	"github.com/dddplayer/core/service"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		os.Exit(-1)
	}

	dot := service.Dot(args[0], args[1], &Repo{
		Tree: radix.NewTree(),
	})

	writeToDisk(dot, strings.ReplaceAll(args[1], "/", "."))
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
