package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"
	"tcp-chatroom-redo/internal/protocol"
)

var (
	clients = map[string]*protocol.Client{}
	mu      sync.Mutex
)

func main() {

	// Create socket, bind, and listen on this address-configuration
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		// Listen for any incoming connections
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}

		// Create a goroutine for every client connection
		go func() {
			defer conn.Close()

			// Loop and receieve messages from our clients
			r := bufio.NewReader(conn)
			for {
				line, err := r.ReadSlice('\n')
				if err != nil { /* handle ErrBufferFull, or disconnect */
				}

				// Trim any trailing \n or \r from our received client message
				m := protocol.ParseMessage(bytes.TrimRight(line, "\r\n"))
				RouteMessage(m.Msg, conn)
			}

		}()
	}
}
