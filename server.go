package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	io.Copy(conn, conn)
	conn.Close()
}

// StartServer starts a node
func StartServer(nodeID int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", nodeID))
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn)
	}
}
