package bc

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
)

func (cli *CLI) getPubKey(privateKey string) {
	curve := elliptic.P256()
	priv_key, _ := hex.DecodeString(privateKey)
	x, y := curve.ScalarBaseMult(priv_key)
	pubKey := append(x.Bytes(), y.Bytes()...)
	fmt.Println(hex.EncodeToString(pubKey))
}

func (cli *CLI) generatePrivKey() {
	private, _ := newKeyPair()
	fmt.Println(hex.EncodeToString(private.D.Bytes()))
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
