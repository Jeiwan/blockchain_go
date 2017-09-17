package main

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

// UTXOSet represents UTXO set
type UTXOSet struct {
	Blockchain *Blockchain
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		bucketName := []byte(utxoBucket)
		b := tx.Bucket(bucketName)

		if b != nil {
			err := tx.DeleteBucket(bucketName)
			if err != nil {
				log.Panic(err)
			}
		}

		_, err := tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindAllUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
}

// GetCount returns the number of transactions in the UTXO set
func (u UTXOSet) GetCount() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// FindUTXO finds UTXO for a public key hash
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}
