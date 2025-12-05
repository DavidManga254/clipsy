package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	db "github.com/DavidManga254/clipsy/lib/DB"
	clipboarditemcontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/clipBoardItemController"
	clip "github.com/DavidManga254/clipsy/lib/clipboard"
	"github.com/DavidManga254/clipsy/lib/logger"
)

func main() {
	logger.Info("Starting up Clipsy oe.....")

	//init db 
	db.InitDB()


	//subscriber to clipboard 
	watcherCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clip.SubscribeToClipboard(watcherCtx)

	clipboarditems, err := clipboarditemcontroller.ListClipboardItems(watcherCtx,db.Database,1,10)

	if err != nil {
		logger.Error("%v",err)
	}

	logger.Info("%v",clipboarditems)
	

	exit()

}


func exit () {
	signalChannel := make(chan os.Signal,1)

	signal.Notify(signalChannel,os.Interrupt,syscall.SIGTERM)
	defer signal.Stop(signalChannel)

	sig:= <- signalChannel

	logger.Info("Exiting program %v",sig)
}