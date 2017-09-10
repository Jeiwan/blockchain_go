package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig []byte
}

// UnlocksOutputWith checks whether the address initiated the transaction
func (in *TXInput) UnlocksOutputWith(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.ScriptSig)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
