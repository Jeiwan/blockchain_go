package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
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

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("Transaction %x:", tx.ID))

	for i, input := range tx.Vin {
		lines = append(lines, fmt.Sprintf("  Input %d:", i))
		lines = append(lines, fmt.Sprintf("    TXID:   %x", input.Txid))
		lines = append(lines, fmt.Sprintf("    Out:    %d", input.Vout))
		lines = append(lines, fmt.Sprintf("    Script: %x", input.ScriptSig))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("  Output %d:", i))
		lines = append(lines, fmt.Sprintf("    Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("    Script: %x", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
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

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, []byte{}})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.ScriptPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTx *Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, vin := range tx.Vin {
		if bytes.Compare(vin.Txid, prevTx.ID) != 0 {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		txCopy.Vin[inID].ScriptSig = prevTx.Vout[vin.Vout].ScriptPubKey
		txCopy.SetID()
		txCopy.Vin[inID].ScriptSig = []byte{}

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].ScriptSig = append(signature, tx.Vin[inID].ScriptSig...)
	}
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTx *Transaction) bool {
	sigLen := 64

	if tx.IsCoinbase() {
		return true
	}

	for _, vin := range tx.Vin {
		if bytes.Compare(vin.Txid, prevTx.ID) != 0 {
			log.Panic("ERROR: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		txCopy.Vin[inID].ScriptSig = prevTx.Vout[vin.Vout].ScriptPubKey
		txCopy.SetID()
		txCopy.Vin[inID].ScriptSig = []byte{}

		signature := vin.ScriptSig[:sigLen]
		pubKey := vin.ScriptSig[sigLen:]

		r := big.Int{}
		s := big.Int{}
		r.SetBytes(signature[:(sigLen / 2)])
		s.SetBytes(signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(pubKey)
		x.SetBytes(pubKey[:(keyLen / 2)])
		y.SetBytes(pubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig []byte
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value        int
	ScriptPubKey []byte
}

// UnlocksOutputWith checks whether the address initiated the transaction
func (in *TXInput) UnlocksOutputWith(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.ScriptSig)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.ScriptPubKey = pubKeyHash
}

// Unlock checks if the output can be used by the owner of the pubkey
func (out *TXOutput) Unlock(pubKeyHash []byte) bool {
	return bytes.Compare(out.ScriptPubKey, pubKeyHash) == 0
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := bc.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
