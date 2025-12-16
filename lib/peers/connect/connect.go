package connectPeers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	db "github.com/DavidManga254/clipsy/lib/DB"
	appmetacontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/appMetaController"
	peerscontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/peerController"
	"github.com/DavidManga254/clipsy/lib/logger"
)

type HandshakeMessage struct {
	DeviceId string
	Message string 
	RecepientDeviceId string
}

func ConnectToPeer(deviceId string,ctx context.Context) (*net.Conn,error) {
	peer, err := peerscontroller.GetPeerById(db.Database,deviceId,ctx)

	if err != nil {
		logger.Error("Failed to retreive peer by id %v",err)

		return nil,err
	}

	myDeviceId,err := appmetacontroller.GetOrCreateDeviceID(ctx, db.Database)

	if err != nil {
		logger.Error("Error retrieving personal device id  %v",err)

		return nil,err
	}

	myDeviceIdint, err := strconv.ParseInt(myDeviceId,0,64)

	if err != nil {
		logger.Error("Failed to parse  My device id string to int %v",err)

		return nil,err
	}

	peerDeviceIdInt, err := strconv.ParseInt(peer.DeviceID,0,64)

	if err != nil {
		logger.Error("Failed to parse  My device id string to int %v",err)

		return nil,err
	}

	logger.Info("My Device id %v",myDeviceId)
	logger.Info("Peer Device id %v",peerDeviceIdInt)

	if myDeviceIdint < peerDeviceIdInt {

		return nil,fmt.Errorf("peer device has higher connect priorityS")
	}
	
	//init connection 
	tcpConn,err := net.DialTimeout("tcp4",peer.IP,2 * time.Second)

	if err!= nil{
		logger.Error("failed to make tcp4 net conn :: %v",err)

		return nil,err
	}

	handshakeMessage := HandshakeMessage{
		DeviceId: myDeviceId,
		Message: "Hello" + peer.DeviceID,
		RecepientDeviceId: peer.DeviceID,
	}

	

	handshakeMessageByte, err := json.Marshal(handshakeMessage)

	if err!= nil{
		logger.Error("failed to convert message to json :: %v",err)

		tcpConn.Close()

		return nil,err
	}

	handshakeMessageByte = append(handshakeMessageByte, '\n')

	//send handshake message 
	var totalBytesSent int = 0 

	for totalBytesSent < len(handshakeMessageByte){
		tcpConn.SetWriteDeadline(time.Now().Add(3 * time.Second))

		sentByteNum,err := tcpConn.Write(handshakeMessageByte[totalBytesSent:])

		if err != nil {
			logger.Info("Error sending handshake message %v",err)

			tcpConn.Close()

			return nil,err
		}

		totalBytesSent += sentByteNum
	}

	//once hand shake completed start reading
	tcpConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(tcpConn)

	handshakeReply, err := reader.ReadBytes('\n')

	if err != nil {
		logger.Info("Error sending handshake message %v",err)

		return nil,err 
	}


	handshakeReply = bytes.TrimSuffix(handshakeReply,[]byte{'\n'})
	handshakeReply = bytes.TrimSuffix(handshakeReply, []byte{'\r'})

	var handshakeData HandshakeMessage

	err = json.Unmarshal(handshakeReply,&handshakeMessage)

	if err != nil {
		logger.Error("Error deoding received json",)
	}

	if handshakeData.DeviceId != peer.DeviceID {
		logger.Error("Received handshake is not desired peer")

		tcpConn.Close()

		return  nil,fmt.Errorf("received handshake is not desired peer")
	}

	if handshakeData.RecepientDeviceId != myDeviceId {
		logger.Error("Handshake reply not addressed to me")

		tcpConn.Close()
		
		return nil,fmt.Errorf("handshake reply not addressed to me")
	}

	return &tcpConn,nil 
}