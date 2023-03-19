package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/goforbroke1006/onix/cmd"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	zap.ReplaceGlobals(logger)

	if err := cmd.ExecuteCmdTree(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
