package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const walletFile = "wallet.dat"

// Wallet ...
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() []byte {
	public := append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...)
	publicSHA256 := sha256.Sum256(public)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	versionedPayload := append([]byte{version}, publicRIPEMD160...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	return address
}

// SaveToFile saves the wallet to a file
func (w Wallet) SaveToFile() {
	var content bytes.Buffer

	gob.Register(w.PrivateKey.Curve)
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

// NewWallet creates and returns a Wallet
func NewWallet() (*Wallet, error) {
	if _, err := os.Stat(walletFile); !os.IsNotExist(err) {
		return nil, errors.New("Wallet already exists")
	}
	private, public := newKeyPair()
	wallet := Wallet{private, public}

	return &wallet, nil
}

func newKeyPair() (ecdsa.PrivateKey, ecdsa.PublicKey) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	return *private, private.PublicKey
}

// Checksum ...
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]
}
