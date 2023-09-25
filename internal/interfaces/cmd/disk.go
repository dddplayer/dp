package cmd

import (
	"fmt"
	"os"
)

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
