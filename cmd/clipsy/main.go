package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	db "github.com/DavidManga254/clipsy/lib/DB"
	appmetacontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/appMetaController"
	clip "github.com/DavidManga254/clipsy/lib/clipboard"
	"github.com/DavidManga254/clipsy/lib/logger"
	connectPeers "github.com/DavidManga254/clipsy/lib/peers/connect"
	peerDiscovery "github.com/DavidManga254/clipsy/lib/peers/discovery"
)

func main() {
	logger.Info("Starting up Clipsy .....")

	watcherCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//init db 
	db.InitDB()

	//init device metadata 
	_,err := appmetacontroller.GetOrCreateDeviceID(watcherCtx,db.Database)
	
	if err != nil {
		logger.Info("failed to init device db metadata %v", err)

		return
	}

	//open udp port 
	udpConn,udpAdrr, err := peerDiscovery.InitUdpPort()

	if err != nil || udpConn == nil {
		logger.Error("error initializing the udp port, closing :: %v",err)

		return
	}
	defer udpConn.Close()


	//open tcp listener and save its details to be communicated during udp handshake 
	tcpListener, err := connectPeers.InitTcpListener()

	if err != nil {
		logger.Info("Failed to initialize tcp port")

		return
	}

	defer tcpListener.Close()

	tcpAddress := tcpListener.Addr().(*net.TCPAddr)

	tcpPort := tcpAddress.Port

	logger.Info("tcpPort :: %v",tcpPort)


	err = appmetacontroller.SetDeviceTcpDetails(watcherCtx,db.Database,tcpPort)

	if err != nil {
		logger.Info("error saving device tcp details")

		return
	}

	//subscriber to clipboard 
	clip.SubscribeToClipboard(watcherCtx)

	// clipboarditems, err := clipboarditemcontroller.ListClipboardItems(watcherCtx,db.Database,1,10)

	// if err != nil {
	// 	logger.Error("%v",err)
	// }

	// logger.Info("%v",clipboarditems)

	go peerDiscovery.DiscoverPeers(udpConn,watcherCtx)
	
	go peerDiscovery.PingPeers(udpConn,udpAdrr,watcherCtx)

	exit()

}


func exit () {
	signalChannel := make(chan os.Signal,1)

	signal.Notify(signalChannel,os.Interrupt,syscall.SIGTERM)
	defer signal.Stop(signalChannel)

	sig:= <- signalChannel

	logger.Info("Exiting program %v",sig)
}