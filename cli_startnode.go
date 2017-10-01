package main

import "fmt"

func (cli *CLI) startNode(nodeID string) {
	fmt.Printf("Starting node %s\n", nodeID)
	StartServer(nodeID)
}
