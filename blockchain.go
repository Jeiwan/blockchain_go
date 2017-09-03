package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbase = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs for an address
func (bc *Blockchain) FindUnspentTransactions(address string) []*Transaction {
	var spentTXs map[string][]int
	var unspentTXs []*Transaction
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txid := string(tx.GetHash())

			for outid, out := range tx.Vout {
				// Was the output spent?
				if spentTXs[txid] != nil {
					for _, spentOut := range spentTXs[txid] {
						if spentOut == outid {
							continue
						}
					}
				}

				if out.Unlock(address) {
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.LockedBy(address) {
						spentTXs[string(in.Txid)] = append(spentTXs[string(in.Txid)], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindUTXOs finds and returns unspend transaction outputs for the address
func (bc *Blockchain) FindUTXOs(address string, amount int) (int, map[string][]int) {
	var unspentOutputs map[string][]int
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

	// TODO: Fix this
	if amount == -1 {
		accumulated = -2
	}

Work:
	for _, tx := range unspentTXs {
		txid := string(tx.GetHash())

		for outid, out := range tx.Vout {
			if out.Unlock(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txid] = append(unspentOutputs[txid], outid)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// Iterator ...
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// Next returns next block starting from the tip
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain(address string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			cbtx := NewCoinbaseTX(address, genesisCoinbase)
			genesis := NewGenesisBlock(cbtx)

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
