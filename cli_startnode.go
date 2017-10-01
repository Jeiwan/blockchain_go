package main

import "fmt"

func (cli *CLI) startNode(nodeID int) {
	fmt.Printf("Starting node %d\n", nodeID)
	StartServer(nodeID)
}
