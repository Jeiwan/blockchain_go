package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const protocol = "tcp"
const dnsNodeID = 3000
const nodeVersion = 1
const commandLength = 12

var nodeAddress string

type verzion struct {
	Version  int
	AddrFrom string
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

func extractCommand(request []byte) []byte {
	return request[:commandLength]
}

func sendVersion(addr string) {
	var payload bytes.Buffer

	enc := gob.NewEncoder(&payload)
	err := enc.Encode(verzion{nodeVersion, nodeAddress})
	if err != nil {
		log.Panic(err)
	}

	request := append(commandToBytes("version"), payload.Bytes()...)

	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	fmt.Printf("%x\n", request)
	_, err = io.Copy(conn, bytes.NewReader(request))
	if err != nil {
		log.Panic(err)
	}
}

func handleConnection(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])

	switch command {
	case "version":
		fmt.Printf("Received %s command", command)
		// send verack
		// send addr
	default:
		fmt.Println("Unknown command received!")
	}

	conn.Close()
}

// StartServer starts a node
func StartServer(nodeID int) {
	nodeAddress = fmt.Sprintf("localhost:%d", nodeID)
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	if nodeID != dnsNodeID {
		sendVersion(fmt.Sprintf("localhost:%d", dnsNodeID))
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn)
	}
}
