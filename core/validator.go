package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type blockValidator struct {
	blockChain *Blockchain
}

func NewBlockValidator(bc *Blockchain) Validator {
	return &blockValidator{blockChain: bc}
}

func (bv *blockValidator) ValidateBlock(block *Block) error {
	if bv.blockChain.HasBlock(block.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", block.Height, block.Hash(BlockHasher))
	}

	return block.Verify()
}
