package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tgragnato/pure/pkg/spam"
)

func main() {

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	httpWorker := make(chan spam.Spam, 1)
	stopped := false

	for !stopped {
		select {

		case <-signalCh:
			stopped = true

		case s := <-httpWorker:
			go s.Call()

		case httpWorker <- spam.MakeSpam():
			time.Sleep(time.Duration(spam.Counter/500) * time.Millisecond)
		}
	}
}
