package main

import (
	"fmt"
	"net"
	"tcp-chatroom-redo/internal/protocol"
)

func SendError(conn net.Conn, errCode protocol.ErrorCode) {
	msg := fmt.Sprintf("ERR %d", errCode)
	err := protocol.Send([]byte(msg), conn)
	if err != nil {
		fmt.Println(err)
	}
}

func isRegistered(conn net.Conn, clients map[string]*protocol.Client) bool {
	for _, c := range clients {
		if c.Conn == conn {
			return true
		}
	}

	return false
}

func handleBroadcastMsg(msg []byte, sender *protocol.Client, clients map[string]*protocol.Client) {

	if len(msg) < protocol.MinMessageLength || len(msg) > protocol.MaxMessageLength {
		SendError(sender.Conn, protocol.InvalidMessage)
		return
	}

	// Otherwise, we broadcast the message to all the clients

	protocol.BroadcastMsg(msg, sender.Conn, clients)

	// Lock for reading
	mu.Lock()
	senderUsername := protocol.NameOf(sender, clients)
	mu.Unlock()

	// Server-log
	fmt.Printf("(%s) -> ALL: %s\n", senderUsername, string(msg))
}

func handleUserRegistration(msg []byte, conn net.Conn, clients map[string]*protocol.Client) {

}

func RouteMessage(readBuf []byte, conn net.Conn) {

	msg := protocol.ParseMessage(readBuf)

	// Handle user-auth; ONLY REG allowed before routing

	if msg.MsgType != protocol.REG && msg.MsgType != protocol.UNKNOWN && !isRegistered(conn, clients) {
		fmt.Println("A client attempted to message before registering.")
		SendError(conn, protocol.InvalidMessage)
	}

	// Route message

	switch msg.MsgType {
	case protocol.REG:
		if isRegistered(conn, clients) {
			SendError(conn, protocol.InvalidMessage)
			return
		}

		handleUserRegistration(readBuf, conn, clients)
		break
	case protocol.MESG:
		protocol.BroadcastMsg(msg.Msg, conn, clients)
		break
	default:
	}

}
