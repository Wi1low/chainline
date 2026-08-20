package core

type Blockchain struct {
	headers   []*Header
	validator Validator
	store     Storage
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
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// Height excludes the genesis block
func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) addBlockWithoutValidator(b *Block) error {
	bc.headers = append(bc.headers, b.Header)

	return bc.store.Put(b)
}
