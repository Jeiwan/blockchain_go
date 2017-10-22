package bc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
	Height        int
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := NewMerkleTree(transactions)

	return mTree.RootNode.Data
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func (b *Block) PrintHTML(detail bool) string {
	var lines []string
	lines = append(lines, fmt.Sprintf("<h2>Block <a href=\"/block/%x\">%x</a> </h2>", b.Hash, b.Hash))
	lines = append(lines, fmt.Sprintf("Height: %d</br>", b.Height))
	lines = append(lines, fmt.Sprintf("Prev. block: <a href=\"/block/%x\">%x</a></br>", b.PrevBlockHash, b.PrevBlockHash))
	lines = append(lines, fmt.Sprintf("Created at : %s</br>", time.Unix(b.Timestamp, 0)))
	pow := NewProofOfWork(b)
	lines = append(lines, fmt.Sprintf("PoW: %s</br></br>", strconv.FormatBool(pow.Validate())))
	if detail {
		for _, tx := range b.Transactions {
			lines = append(lines, tx.PrintHTML())
		}
	}
	lines = append(lines, fmt.Sprintf("</br></br>"))
	return strings.Join(lines, "\n")
}
