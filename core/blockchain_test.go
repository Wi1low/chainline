package core

import (
	"testing"

	"github.com/Wi1low/chainline/types"
	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0), NewMemoryStore())
	assert.Nil(t, err)

	return bc
}
func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher.Hash(prevHeader)
}

// TestAddBlock 是一个测试函数，用于测试区块链添加区块的功能
func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	length := 1000
	for i := 0; i < length; i++ {
		b := randomBlockWithSignatureAndPrevBlockHash(
			t,
			uint32(i+1),
			getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(b))
	}
	assert.Equal(t, bc.Height(), uint32(length))
	assert.Equal(t, len(bc.headers), length+1)

	// Add block with invalid height
	assert.NotNil(t, bc.AddBlock(randomBlock(77)))

	// Add block with invalid height
	assert.NotNil(t, bc.AddBlock(randomBlock(3000)))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	length := 1000

	for i := 0; i < length; i++ {
		b := randomBlockWithSignatureAndPrevBlockHash(
			t,
			uint32(i+1),
			getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(b))

		header, err := bc.GetHeader(b.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, b.Header)

	}

}
func TestBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.False(t, bc.HasBlock(1))
	assert.True(t, bc.HasBlock(0))
}
