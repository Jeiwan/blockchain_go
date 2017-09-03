package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const subsidy = 10

// Transaction represents a Bitcoin transaction
type Transaction struct {
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Txid == -1 && tx.Vin[0].Vout == -1
}

// GetHash hashes the transaction and returns the hash
func (tx Transaction) GetHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())

	return hash[:]
}

// TXInput represents a transaction input
type TXInput struct {
	Txid      int
	Vout      int
	ScriptSig string
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// Unlock checks if the output can be unlocked with the provided data
func (out *TXOutput) Unlock(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = "Coinbase"
	}

	txin := TXInput{-1, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{[]TXInput{txin}, []TXOutput{txout}}

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, value int) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := s.findUnspentOutputs(from, value)

	if acc < value {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		for _, out := range outs {
			input := TXInput{txid, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{value, to})
	if acc > value {
		outputs = append(outputs, TXOutput{acc - value, from}) // a change
	}

	tx := Transaction{inputs, outputs}

	return &tx
}
