package main

import "fmt"

func (cli *CLI) reindexUTXO() {
	bc := NewBlockchain()
	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	count := UTXOSet.GetCount()
	fmt.Printf("Done! There are %d transactions in the UTXO set.", count)
}
