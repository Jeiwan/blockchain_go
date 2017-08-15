package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Block keeps block headers
type Block struct {
	Timestamp int64
	Data      []byte
	PrevBlock []byte
	hash      []byte
}

// SetHash calculates and sets block hash
func (b *Block) SetHash() {
	h := sha256.New()
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	var data []byte
	data = append(data, b.PrevBlock...)
	data = append(data, b.Data...)
	data = append(data, timestamp...)

	_, err := h.Write(data)
	if err != nil {
		log.Panic(err)
	}
	b.hash = h.Sum(nil)
}

// NewBlock creates and returns Block
func NewBlock(data string, prevBlock []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlock, []byte("")}
	block.SetHash()
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte("0"))
}

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	blocks []*Block
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := &Block{time.Now().Unix(), []byte(data), prevBlock.hash, []byte("")}
	newBlock.SetHash()
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("%s\n", block.Data)
		fmt.Printf("%x\n", block.hash)
	}
}
