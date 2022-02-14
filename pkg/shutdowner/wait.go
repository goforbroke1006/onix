package shutdowner

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\r- Ctrl+C pressed in Terminal")
}
