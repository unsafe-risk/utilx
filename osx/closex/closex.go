package closex

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type interrupt struct {
	C chan struct{}
}

func HandleKillSig(handler func(), logger *log.Logger) interrupt {
	intr := interrupt{
		C: make(chan struct{}),
	}
	sigChannel := make(chan os.Signal, 1)

	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for signal := range sigChannel {
			logger.Printf("Receive signal %s, Shutting down...", signal)
			handler()
			close(intr.C)
		}
	}()
	return intr
}
