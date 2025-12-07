package peers

/*
two features

1.discovery
2. Messaging

api expectation
start - discovery and peer search
stop - stop module
send - data send features
peers - list of all known peers

will use udp broadcast for multicast based discovery
each node open a fixed port
periedically send a broadcast
listen on the same port for discovery messages
when receives a response, register that as a peer

discovery message format:
-message type
-current device id
-port
-metadata : device name etc....

data messages
-type
-from - sender id
-to - recepient
-payload TODO decide

Peer

ID – unique string.

Addr – IP:port for TCP connection.

Conn – the active TCP connection to that peer (if any).

LastSeen – timestamp of last discovery message.

Manager / Node

ID – your own node ID.

Peers – map of peer ID to Peer struct.

UDPConn – UDP socket for discovery.

TCPListener – TCP listener for incoming data connections.

Incoming – channel or callback for received data.

Quit / Context – for stopping goroutines cleanly.

You will run different goroutines for:

Discovery listener loop (read UDP, handle messages).

Discovery broadcaster loop (send periodic hello).

TCP accept loop (accept incoming connections).

TCP per-connection loops (read incoming messages from that peer).

*/