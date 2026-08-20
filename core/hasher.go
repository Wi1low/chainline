package core

import (
	"crypto/sha256"

	"github.com/Wi1low/chainline/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

// implements Hashers

var BlockHasher Hasher[*Block] = &blockHasher{}

type blockHasher struct{}

func (blockHasher) Hash(block *Block) types.Hash {
	sha := sha256.Sum256(block.HeaderData())
	return types.Hash(sha)
}
