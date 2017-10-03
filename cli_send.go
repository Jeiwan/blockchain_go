package main

import (
	"fmt"
	"log"
)

func (cli *CLI) send(from, to string, amount int, nodeID string) {
	if !ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := NewBlockchain(nodeID)
	UTXOSet := UTXOSet{bc}
	defer bc.db.Close()

	wallets, err := NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	tx := NewUTXOTransaction(&wallet, to, amount, &UTXOSet)
	cbTx := NewCoinbaseTX(from, "")
	// txs := []*Transaction{cbTx, tx}

	// var txHashes [][]byte
	// txHashes = append(txHashes, tx.Hash())
	// txHashes = append(txHashes, cbTx.Hash())

	sendTx(knownNodes[0], tx)
	sendTx(knownNodes[0], cbTx)

	// newBlock := bc.MineBlock(txs)
	// UTXOSet.Update(newBlock)
	fmt.Println("Success!")
}
