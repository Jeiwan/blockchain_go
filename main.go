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

func main() {
	gb := NewGenesisBlock()
	b1 := NewBlock("Send 1 BTC to Ivan", gb.hash)

	fmt.Printf("%s\n", gb.Data)
	fmt.Printf("%x\n", gb.hash)
	fmt.Printf("%s\n", b1.Data)
	fmt.Printf("%x\n", b1.hash)
}
