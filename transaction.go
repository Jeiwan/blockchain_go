package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction represents a Bitcoin transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID sets ID of a transaction
func (tx Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, value int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindUTXOs(from, value)

	if acc < value {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		for _, out := range outs {
			txID, err := hex.DecodeString(txid)
			if err != nil {
				log.Panic(err)
			}

			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{value, to})
	if acc > value {
		outputs = append(outputs, TXOutput{acc - value, from}) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
