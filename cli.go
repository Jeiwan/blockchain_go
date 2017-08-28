package main

import (
	"fmt"
	"os"
	"strconv"
)

// CLI responsible for processing command line arguments
type CLI struct {
	bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addBlock BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printChain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock() {
	if len(os.Args) < 3 {
		fmt.Println("Error: BLOCK_DATA is not specified.")
		os.Exit(1)
	}
	cli.bc.AddBlock(os.Args[2])
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

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

// ProcessArgs parses command line arguments and processes commands
func (cli *CLI) ProcessArgs() {
	cli.validateArgs()

	switch os.Args[1] {
	case "addBlock":
		cli.addBlock()
	case "printChain":
		cli.printChain()
	default:
		cli.printUsage()
		os.Exit(1)
	}
}
