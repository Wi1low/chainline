package core

import "sync"

type Storage interface {
	// Get() (*Block, error)
	Put(*Block) error
}

type MemoryStorage struct {
	rwlock sync.RWMutex
}

func NewMemoryStore() Storage {
	return &MemoryStorage{}
}
func (ms *MemoryStorage) Put(b *Block) error {

	return nil
}
