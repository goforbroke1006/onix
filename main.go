package main

import (
	"fmt"
	"os"

	"github.com/goforbroke1006/onix/cmd"
)

func main() {
	if err := cmd.ExecuteCmdTree(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
