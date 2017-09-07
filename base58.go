package main

import (
	"math/big"
)

var alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode encodes a byte array to Base58
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0)
	x.SetBytes(input)

	base := big.NewInt(int64(len(alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	for c := range input {
		if c == 0x00 {
			result = append([]byte{alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}
