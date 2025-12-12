package peerDiscovery

import (
	"context"
	"encoding/json"
	"net"
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

/*
use udp socket to discover peers
*/
//TODO return statements will stop the loop check how to handle to allow for more reading 
func DiscoverPeers(ctx context.Context) error {
	udpAdrr, err := net.ResolveUDPAddr("udp4",":40000")

	if err != nil {
		logger.Error("Error allocating address :: %v",err)

		return err
	}

	udpConn, err := net.ListenUDP("udp4", udpAdrr)
	if err != nil {
		logger.Error("Error initializing listen udp :: %v", err)
		return err
	}
	defer udpConn.Close()

	var udpPayload []byte = make([]byte, 2048)

	for {
		select {
			case <-ctx.Done():
				return ctx.Err()
			default:
		}

		err = udpConn.SetReadDeadline(time.Now().Add(3 * time.Second))

		if err != nil {
			logger.Error("Error setting read deadline :: %v", err)

			return err
		}

		logger.Info("reading")


		byteNumCopied, address, err := udpConn.ReadFromUDP(udpPayload)

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

		err = peerscontroller.New(db.Database,peerDbData,ctx)

		if err != nil {
			logger.Error("error adding resolved peer to DB :: %v", err)
			continue
		}
	}
}

func PingPeers(ctx context.Context) error{
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

	udpAdrr, err := net.ResolveUDPAddr("udp4","127.0.0.1:40000")

	if err != nil {
		logger.Error("Error allocating address :: %v",err)

		return err
	}

	udpConn,err := net.DialUDP("udp4",nil,udpAdrr)

	if err != nil {
		logger.Error("Error dialing up %v",err)

		return err
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
				_,err := udpConn.Write(marshalMessage)

				if err != nil {
					logger.Error("error sending ping message %v",err)

					continue
				}
			}
		}
	}
}