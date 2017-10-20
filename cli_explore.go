package main

import (
	"encoding/hex"
	"fmt"
)

func (cli *CLI) generateKey() {
	private, public := newKeyPair()
	fmt.Println("Public Key:")
	fmt.Println(hex.EncodeToString(public))
	priKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	priKey = append(priKey, private.D.Bytes()...)
	fmt.Println("Private Key:")
	fmt.Printf("%d%s\n", private.Curve, hex.EncodeToString(priKey))
}

func (cli *CLI) getAddress(pubKey string) {
	public, _ := hex.DecodeString(pubKey)

	pubKeyHash := HashPubKey(public)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	fullPayload := append(versionedPayload, checksum(versionedPayload)...)

	fmt.Println()
	fmt.Printf("PubKey     : %s\n", pubKey)
	fmt.Printf("PubKeyHash : %x\n", pubKeyHash)
	fmt.Printf("Address    : %s\n", Base58Encode(fullPayload))
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
