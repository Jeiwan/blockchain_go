package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	sigLen := 64
	pubKey := in.ScriptSig[sigLen:]
	lockingHash := HashPubKey(pubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
