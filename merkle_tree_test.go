package main

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMerkleNode(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}

	// Level 1

	n1 := NewMerkleNode(nil, nil, data[0])
	n2 := NewMerkleNode(nil, nil, data[1])
	n3 := NewMerkleNode(nil, nil, data[2])
	n4 := NewMerkleNode(nil, nil, data[2])

	// Level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// Level 3
	n7 := NewMerkleNode(n5, n6, nil)

	assert.Equal(
		t,
		"64b04b718d8b7c5b6fd17f7ec221945c034cfce3be4118da33244966150c4bd4",
		hex.EncodeToString(n5.Data),
		"Level 1 hash 1 is correct",
	)
	assert.Equal(
		t,
		"08bd0d1426f87a78bfc2f0b13eccdf6f5b58dac6b37a7b9441c1a2fab415d76c",
		hex.EncodeToString(n6.Data),
		"Level 1 hash 2 is correct",
	)
	assert.Equal(
		t,
		"4e3e44e55926330ab6c31892f980f8bfd1a6e910ff1ebc3f778211377f35227e",
		hex.EncodeToString(n7.Data),
		"Root hash is correct",
	)
}

func TestNewMerkleTree(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
	}
	// Level 1
	n1 := NewMerkleNode(nil, nil, data[0])
	n2 := NewMerkleNode(nil, nil, data[1])
	n3 := NewMerkleNode(nil, nil, data[2])
	n4 := NewMerkleNode(nil, nil, data[2])

	// Level 2
	n5 := NewMerkleNode(n1, n2, nil)
	n6 := NewMerkleNode(n3, n4, nil)

	// Level 3
	n7 := NewMerkleNode(n5, n6, nil)

	rootHash := fmt.Sprintf("%x", n7.Data)
	mTree := NewMerkleTree(data)

	assert.Equal(t, rootHash, fmt.Sprintf("%x", mTree.RootNode.Data), "Merkle tree root hash is correct")
}

func TestNewMerkleTree8(t *testing.T) {
	data := [][]byte{
		[]byte("node1"),
		[]byte("node2"),
		[]byte("node3"),
		[]byte("node4"),
		[]byte("node5"),
		[]byte("node6"),
	}
	// Level 0
	n01 := NewMerkleNode(nil, nil, data[0])
	n02 := NewMerkleNode(nil, nil, data[1])
	n03 := NewMerkleNode(nil, nil, data[2])
	n04 := NewMerkleNode(nil, nil, data[3])
	n05 := NewMerkleNode(nil, nil, data[4])
	n06 := NewMerkleNode(nil, nil, data[5])
	n07 := NewMerkleNode(nil, nil, data[5])
	n08 := NewMerkleNode(nil, nil, data[5])

	// Level 1
	n11 := NewMerkleNode(n01, n02, nil)
	n12 := NewMerkleNode(n03, n04, nil)
	n13 := NewMerkleNode(n05, n06, nil)
	n14 := NewMerkleNode(n07, n08, nil)

	// Level 2
	n21 := NewMerkleNode(n11, n12, nil)
	n22 := NewMerkleNode(n13, n14, nil)

	// Level 3
	n31 := NewMerkleNode(n21, n22, nil)

	rootHash := fmt.Sprintf("%x", n31.Data)
	mTree := NewMerkleTree(data)

	assert.Equal(t, rootHash, fmt.Sprintf("%x", mTree.RootNode.Data), "Merkle tree root hash is correct")
}
