package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	bc := NewBlockchain()
	defer bc.db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Wrong!")
		os.Exit(1)
	}

	if os.Args[1] == "addBlock" {
		bc.AddBlock(os.Args[2])
		fmt.Println("Success!")
	}

	if os.Args[1] == "printChain" {
		bci := bc.Iterator()

		for {
			block := bci.Next()

			fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
			fmt.Printf("Data: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			pow := NewProofOfWork(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()

			if len(block.PrevBlockHash) == 0 {
				break
			}
		}
	}
}
