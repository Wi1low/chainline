package core

import (
	"crypto/sha256"

	"github.com/Wi1low/chainline/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

// implements Hashers

var BlockHasher Hasher[*Header] = &blockHasher{}
var TxHasher Hasher[*Transaction] = &txHasher{}

type blockHasher struct{}

func (blockHasher) Hash(header *Header) types.Hash {
	sha := sha256.Sum256(header.Bytes())
	return types.Hash(sha)
}

type txHasher struct{}

func (txHasher) Hash(tx *Transaction) types.Hash {
	sha := sha256.Sum256(tx.Data)
	return types.Hash(sha)
}
