package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/Wi1low/chainline/crypto"
	"github.com/Wi1low/chainline/types"
	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	txs := []Transaction{
		{
			Data: []byte("test block1"),
		},
		{
			Data: []byte("test block2"),
		},
	}

	return NewBlock(header, txs)
}
func randomBlockWithSignature(t *testing.T, height uint32) *Block {
	privkey := crypto.GeneratePrivateKey()

	b := randomBlock(height)
	assert.Nil(t, b.Sign(privkey))

	return b
}

func TestHashBlock(t *testing.T) {
	b := randomBlock(0)
	fmt.Printf("b.Hash(BlockHasher): %v\n", b.Hash(BlockHasher))
}

func TestSignBlock(t *testing.T) {
	privkey := crypto.GeneratePrivateKey()

	b := randomBlock(0)
	assert.Nil(t, b.Sign(privkey))
	assert.NotNil(t, b.Signature)
}
func TestVerifyBlock(t *testing.T) {
	privkey := crypto.GeneratePrivateKey()
	// original
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privkey))

	assert.Nil(t, b.Verify())

	// modify the public key
	otherPrikey := crypto.GeneratePrivateKey()
	b.Validator = otherPrikey.PublicKey()
	assert.NotNil(t, b.Verify())

	// modify the block infomation
	fmt.Printf("before %+v\n", b.Header)
	b.Height = 200
	assert.NotNil(t, b.Verify())
	fmt.Printf("after %+v\n", b.Header)

	// 只有原来的秘钥才能验证数据
	// original
	b.Validator = privkey.PublicKey()
	assert.Nil(t, b.Verify())
}
