package main

import (
	"encoding/hex"
	"fmt"
)

func (cli *CLI) generateKey() {
	_, public := newKeyPair()
	fmt.Println("Public Key:")
	fmt.Println(hex.EncodeToString(public))
}

func (cli *CLI) getAddress(pubKey string) {
	public, _ := hex.DecodeString(pubKey)

	pubKeyHash := HashPubKey(public)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)
	fmt.Println()
	fmt.Printf("PubKey     : %s\n", pubKey)
	fmt.Printf("PubKeyHash : %x\n", pubKeyHash)
	fmt.Printf("Address    : %s\n", address)
}

func (cli *CLI) getPubKeyHash(address string) {
	pubKeyHash := Base58Decode([]byte(address))
	fmt.Printf("%x\n", pubKeyHash[1:len(pubKeyHash)-4])
}

func (cli *CLI) validateAddr(address string) {
	fmt.Printf("Address: %s\n", address)
	if !ValidateAddress(address) {
		fmt.Println("Not valid!")
	} else {
		fmt.Println("Valid!")
	}
}
