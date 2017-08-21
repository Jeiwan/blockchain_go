package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

var (
	target   big.Int
	maxNonce = math.MaxInt64
)

const targetBits = 12

func setTarget() {
	targetBytes := make([]byte, 32)
	numOfZeros := targetBits / 4
	targetBytes[numOfZeros-1] = 1
	target.SetBytes(targetBytes)
}

func prepareData(block *Block, nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			intToHex(block.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// Prove ...
func Prove(block *Block) (int, []byte) {
	setTarget()
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", block.Data)
	for nonce < maxNonce {
		data := prepareData(block, nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(&target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// ConfirmProof ..
func ConfirmProof(block *Block, nonce int) bool {
	setTarget()
	var hashInt big.Int

	data := prepareData(block, nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	confirmation := hashInt.Cmp(&target) == -1

	return confirmation
}
