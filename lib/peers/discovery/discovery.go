package peerDiscovery

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"time"

	db "github.com/DavidManga254/clipsy/lib/DB"
	appmetacontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/appMetaController"
	peerscontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/peerController"
	"github.com/DavidManga254/clipsy/lib/logger"
)

type MessageType string

var (
	DiscoverMessage MessageType = "DISCOVER"
	DataMessage     MessageType = "DATA"
)

type DiscoverPeer struct {
	MesssageType MessageType
	DeviceID     string
	Port         string
	Metadata     struct {
		DeviceName string
	}
}

var portList []string = []string{
	"55000",
	"56000",
	"60000",
	"61000",
	"62000",
	"65000",

}


func InitUdpPort() (*net.UDPConn, *net.UDPAddr ,error) {

	var lastErr error

	var UdpConn *net.UDPConn

	var UpdArr *net.UDPAddr


	for _, port := range portList {
		udpAdrr, err := net.ResolveUDPAddr("udp4",":"+port)

		if err != nil {
			logger.Error("Error allocating addressmfor port %v :: %v",port,err)

			lastErr = err

			continue
		}

		udpConn, err := net.ListenUDP("udp4", udpAdrr)

		if err != nil {
			logger.Error("Error initializing listen udp :: %v", err)

			lastErr = err

			continue
		}

		// file, err := udpConn.File()

		// if err != nil {
		// 	logger.Error("error getting file descriptor for port %v",err)

		// 	udpConn.Close()

		// 	continue
		// }

    	// defer file.Close()

		// //TODO handle casting depending on operating system
		// fd := syscall.Handle(file.Fd())

		// err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)

		// if err != nil {
		// 	udpConn.Close()

		// 	logger.Error("failed to enable broadcast: %v", err)

		// 	continue
		// }

		// lc := net.ListenConfig{
        //     Control: func(network, address string, c syscall.RawConn) error {
        //         var opErr error
        //         err := c.Control(func(fd uintptr) {
        //             // Enable broadcast - works on Windows
        //             opErr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
        //         })
        //         if err != nil {
        //             return err
        //         }
        //         return opErr
        //     },
        // }

        // Listen using the configured listener
        // conn, err := lc.ListenPacket(context.Background(), "udp4", udpAdrr.String())
        // if err != nil {
        //     logger.Error("Error initializing listen udp :: %v", err)
        //     lastErr = err
        //     continue
        // }

        // udpConn := conn.(*net.UDPConn)
        

		lastErr = err

		logger.Info("udp addressed opened at port %v",port)

		UdpConn = udpConn
		UpdArr = udpAdrr

		break
	}

	return UdpConn, UpdArr, lastErr
}

/*
use udp socket to discover peers
*/
//TODO return statements will stop the loop check how to handle to allow for more reading 
func DiscoverPeers(UdpConn *net.UDPConn, ctx context.Context) error {
	var udpPayload []byte = make([]byte, 2048)

	for {
		select {
			case <-ctx.Done():
				return ctx.Err()
			default:
		}

		err := UdpConn.SetReadDeadline(time.Now().Add(3 * time.Second))

		if err != nil {
			logger.Error("Error setting read deadline :: %v", err)

			return err
		}

		logger.Info("reading")


		byteNumCopied, address, err := UdpConn.ReadFromUDP(udpPayload)

		if err != nil {
			// treat timeout as "no data, try again"
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// check context before next read
				if ctx.Err() != nil {
					return ctx.Err()
				}

				continue
			}

			logger.Error("Error reading from udp :: %v", err)
			return err
		}

		if byteNumCopied == 0 || address == nil {
			continue
		}

		payloadJson := udpPayload[:byteNumCopied]

		var payload DiscoverPeer

		err = json.Unmarshal(payloadJson, &payload)

		if err != nil {
			logger.Error("error decoding received payload :: %v", err)
			continue
		}

		deviceId,err := appmetacontroller.GetOrCreateDeviceID(ctx,db.Database)

		if err != nil {
			logger.Error("error getting device id:: %v", err)
			continue
		}

		if deviceId == payload.DeviceID {
			continue
		}

		var peerDbData peerscontroller.Peer = peerscontroller.Peer{
			DeviceID: payload.DeviceID,
			IP: address.IP.String(),
			Port: address.Port,
			DeviceName: payload.Metadata.DeviceName,
			ProtocolID: "",
			Status: peerscontroller.PeerStatusDiscovered,
			Metadata: make(map[string]string),
			CreatedAt: time.Now(),
			LastSeen: time.Now(),
			Is_trusted:0,
		}

		logger.Info("Logged data %v",peerDbData)

		err = peerscontroller.New(db.Database,peerDbData,ctx)

		if err != nil {
			logger.Error("error adding resolved peer to DB :: %v", err)
			continue
		}
	}
}

func PingPeers(udpConn *net.UDPConn,udpAdrr *net.UDPAddr,ctx context.Context) error{
	deviceId,err := appmetacontroller.GetOrCreateDeviceID(ctx,db.Database)

	if err != nil {
		logger.Error("error getting device id:: %v", err)
		
		return err
	}

	deviceInfo := DiscoverPeer{
		MesssageType: DiscoverMessage,
		DeviceID: deviceId,
		Port: "40000",
		Metadata: struct{DeviceName string}{DeviceName: "Dave"},
	}

	marshalMessage, err := json.Marshal(deviceInfo)

	if err != nil {
		logger.Error("Error marshalling the deviceinfo to json %v",err)

		return err
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		logger.Info("sending")
		select{
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:{

				for _, port := range portList{
					portInt,err := strconv.ParseInt(port,10,0)

					if err != nil {
						logger.Error("Failed to convert port string to int  %v",err)

						continue
					}

					portAddress := &net.UDPAddr{
						IP:   net.IPv4bcast,
						Port: int(portInt),
					}
					_,err = udpConn.WriteToUDP(marshalMessage,portAddress)

					if err != nil {
						logger.Error("error sending ping message %v",err)

						continue
					}
				}
				
			}
		}
	}
}