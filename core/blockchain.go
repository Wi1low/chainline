package core

import (
	"fmt"
	"log/slog"
	"sync"
)

type Blockchain struct {
	headers   []*Header
	validator Validator
	store     Storage

	rwlock sync.RWMutex
}

// NewBlockchain creates a new blockchain with the given genesis block
func NewBlockchain(genesis *Block, store Storage) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   store,
	}
	bc.validator = NewBlockValidator(bc)

	err := bc.addBlockWithoutValidator(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidato(validator Validator) {
	bc.validator = validator
}

func (bc *Blockchain) AddBlock(b *Block) error {

	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidator(b)
}
func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {

	if !bc.HasBlock(height) {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.rwlock.RLock()
	defer bc.rwlock.RUnlock()

	return bc.headers[height], nil
}
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// Height excludes the genesis block
func (bc *Blockchain) Height() uint32 {
	bc.rwlock.RLock()
	defer bc.rwlock.RUnlock()
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlockWithoutValidator(b *Block) error {
	bc.rwlock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.rwlock.Unlock()

	slog.Info("adding new block",
		slog.Any("block", b.Header),
	)

	return bc.store.Put(b)
}
