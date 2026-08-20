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
	// check if the block already exists in the chain
	if bv.blockChain.HasBlock(block.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", block.Height, block.Hash(BlockHasher))
	}

	// the height must be one by one, in case invalid height is passed
	if block.Height != bv.blockChain.Height()+1 {
		return fmt.Errorf("block (%s) too high", block.Hash(BlockHasher))
	}

	// prevHeader
	prevHeader, err := bv.blockChain.GetHeader(block.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher.Hash(prevHeader)
	if hash != block.PrevBlockHash {
		return fmt.Errorf("block (%s) prevHash (%s) does not match", block.Hash(BlockHasher), block.PrevBlockHash)
	}
	// finally, check the block itself
	return block.Verify()
}
