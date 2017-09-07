package main

import (
	"bytes"
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

// Base58Decode decodes Base58 data
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for c := range input {
		if c == 0x00 {
			zeroBytes++
		}
	}

	address := input[zeroBytes:]
	for _, b := range address {
		charIndex := bytes.IndexByte(alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	raw := result.Bytes()
	raw = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), raw...)

	return raw
}
