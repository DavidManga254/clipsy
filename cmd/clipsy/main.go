package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	clip "github.com/DavidManga254/clipsy/lib/clipboard"
	"github.com/DavidManga254/clipsy/lib/logger"
)

func main() {
	logger.Info("Starting up Clipsy oe.....")


	//subscriber to clipboard 
	watcherCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clip.SubscribeToClipboard(watcherCtx)

	exit()

}


func exit () {
	signalChannel := make(chan os.Signal,1)

	signal.Notify(signalChannel,os.Interrupt,syscall.SIGTERM)
	defer signal.Stop(signalChannel)

	sig:= <- signalChannel

	logger.Info("Exiting program %v",sig)
}