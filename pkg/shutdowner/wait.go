package shutdowner

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	//fmt.Println("\r- Ctrl+C pressed in Terminal")
}
