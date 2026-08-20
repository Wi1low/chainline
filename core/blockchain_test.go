package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0), NewMemoryStore())
	assert.Nil(t, err)

	return bc
}

func TestAddBlock(t *testing.T) {
	bc := randomBlockchainWithGenesis(t)
	length := 1000
	for i := 0; i < length; i++ {
		b := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(b))
	}
	assert.Equal(t, bc.Height(), uint32(length))
	assert.Equal(t, len(bc.headers), length+1)

	// Add block with invalid height
	assert.NotNil(t, bc.AddBlock(randomBlock(77)))
}

func TestBlockchain(t *testing.T) {
	bc := randomBlockchainWithGenesis(t)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := randomBlockchainWithGenesis(t)

	assert.False(t, bc.HasBlock(1))
	assert.True(t, bc.HasBlock(0))
}
