package main

import "log"

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
func NewCoinbaseTX(to string) *Transaction {
	txin := TXInput{-1, -1, "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"}
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
