package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	db "github.com/DavidManga254/clipsy/lib/DB"
	clip "github.com/DavidManga254/clipsy/lib/clipboard"
	"github.com/DavidManga254/clipsy/lib/logger"
	peerDiscovery "github.com/DavidManga254/clipsy/lib/peers/discovery"
)

func main() {
	logger.Info("Starting up Clipsy .....")

	//init db 
	db.InitDB()


	//subscriber to clipboard 
	watcherCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clip.SubscribeToClipboard(watcherCtx)

	// clipboarditems, err := clipboarditemcontroller.ListClipboardItems(watcherCtx,db.Database,1,10)

	// if err != nil {
	// 	logger.Error("%v",err)
	// }

	// logger.Info("%v",clipboarditems)

	go peerDiscovery.DiscoverPeers(watcherCtx)
	
	go peerDiscovery.PingPeers(watcherCtx)

	exit()

}


func exit () {
	signalChannel := make(chan os.Signal,1)

	signal.Notify(signalChannel,os.Interrupt,syscall.SIGTERM)
	defer signal.Stop(signalChannel)

	sig:= <- signalChannel

	logger.Info("Exiting program %v",sig)
}