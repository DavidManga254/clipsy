package peerDiscovery

import (
	"context"
	"encoding/json"
	"net"
	"time"

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

		logger.Info("here is received address %v", address)
		logger.Info("Here is received data :: %v", payload.DeviceID)
	}
}