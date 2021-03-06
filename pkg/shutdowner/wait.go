package shutdowner

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForShutdown locks goroutine and wait for Ctrl+C.
func WaitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
