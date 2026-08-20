package core

type Storage interface {
	// Get() (*Block, error)
	Put(*Block) error
}

type MemoryStorage struct{}

func NewMemoryStore() Storage {
	return &MemoryStorage{}
}
func (ms *MemoryStorage) Put(b *Block) error {

	return nil
}
