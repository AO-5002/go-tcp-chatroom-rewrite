package protocol

import (
	"fmt"
	"io"
	"net"
)

type MSGTYPE int

const (
	REG MSGTYPE = iota
	ACK
	MESG
	PMSG
	ERR
	EXIT
	UNKNOWN
)

type MSGProtocol struct {
	MsgType MSGTYPE
	Msg     []byte
}

type Client struct {
	Conn     net.Conn
	SendChan chan []byte
}

// NameOf Reverse-lookup
func NameOf(client *Client, clients map[string]*Client) string {
	for name, c := range clients {
		if c == client {
			return name
		}
	}
	return ""
}

func Recv(msg []byte, conn net.Conn) (int, error) {
	n, err := conn.Read(msg)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Connection closed.:", err)
			return 0, nil
		}

		return 0, err
	}

	return n, nil
}

func BroadcastMsg(msg []byte, originConn net.Conn, clients map[string]*Client) {
	for _, c := range clients {
		if originConn != c.Conn {
			err := Send(msg, c.Conn)
			if err != nil {
				return
			}
		}
	}
}

func Send(msg []byte, conn net.Conn) error {
	_, err := conn.Write(msg)
	return err
}

func getMsgType(msg []byte) MSGTYPE {
	switch string(msg) {
	case "REG":
		return REG
	case "ACK":
		return ACK
	case "MESG":
		return MESG
	case "PMSG":
		return PMSG
	case "ERR":
		return ERR
	case "EXIT":
		return EXIT
	default:
		return UNKNOWN
	}
}

func ParseMessage(msg []byte) MSGProtocol {
	// Parse the type and content from the message.
	i := 0
	for ; i < 4 && i < len(msg); i++ {
		if msg[i] == ' ' {
			break
		}
	}
	msgType := msg[:i]

	// Type must be followed by a space or the end of the string.
	if i < len(msg) && msg[i] != ' ' {
		return MSGProtocol{UNKNOWN, msg[i:]}
	}

	if i < len(msg) {
		i++
	}

	return MSGProtocol{getMsgType(msgType), msg[i:]}
}
