package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase58(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, public := newKeyPair()
		pubKeyHash := HashPubKey(public)

		versionedPayload := append([]byte{version}, pubKeyHash...)
		checksum := checksum(versionedPayload)

		fullPayload := append(versionedPayload, checksum...)
		address := Base58Encode(fullPayload)

		assert.Equal(
			t,
			ValidateAddress(string(address[:])),
			true,
			"Address: %s is invalid", address,
		)
	}
}
