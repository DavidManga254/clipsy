package main

import (
	"context"
	"fmt"
	"net"

	clip "github.com/DavidManga254/clipsy/lib/clipboard"
	db "github.com/DavidManga254/clipsy/lib/db"
	appmetacontroller "github.com/DavidManga254/clipsy/lib/db/controllers/appMetaController"
	"github.com/DavidManga254/clipsy/lib/logger"
	connectPeers "github.com/DavidManga254/clipsy/lib/peers/connect"
	peerDiscovery "github.com/DavidManga254/clipsy/lib/peers/discovery"
)

// App struct
type App struct {
	ctx context.Context
	udpConn *net.UDPConn
	tcpListener net.Listener

}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}


func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	//init db 
	db.InitDB()

	//init device metadata 
	_,err := appmetacontroller.GetOrCreateDeviceID(ctx,db.Database)
	
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
	a.udpConn = udpConn


	//open tcp listener and save its details to be communicated during udp handshake 
	tcpListener, err := connectPeers.InitTcpListener()

	if err != nil {
		logger.Info("Failed to initialize tcp port")

		return
	}

	a.tcpListener = tcpListener


	tcpAddress := tcpListener.Addr().(*net.TCPAddr)

	tcpPort := tcpAddress.Port

	logger.Info("tcpPort :: %v",tcpPort)

	err = appmetacontroller.SetDeviceTcpDetails(ctx,db.Database,tcpPort)

	if err != nil {
		logger.Info("error saving device tcp details")

		return
	}

	//subscriber to clipboard 
	clip.SubscribeToClipboard(ctx)

	go peerDiscovery.DiscoverPeers(udpConn,ctx)
	
	go peerDiscovery.PingPeers(udpConn,udpAdrr,ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Shutdown(ctx context.Context){
	//close udp port
	a.udpConn.Close()

	//close  tcp port 
	a.tcpListener.Close()

}
