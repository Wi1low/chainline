package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/Wi1low/chainline/types"
	"github.com/stretchr/testify/assert"
)

func TestHeaderDecodeEncode(t *testing.T) {
	h := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Timestamp:     time.Now().UnixNano(),
		Height:        39,
		Nonce:         123,
	}

	buf := &bytes.Buffer{}

	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}
func TestBlockDecodeEncode(t *testing.T) {
	b := &Block{
		Header: &Header{
			Version:       1,
			PrevBlockHash: types.RandomHash(),
			Timestamp:     time.Now().UnixNano(),
			Height:        39,
			Nonce:         123,
		},
		Transactions: nil,
	}
	buf := &bytes.Buffer{}

	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
	fmt.Printf("%+v\n", bDecode)
}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: &Header{
			Version:       1,
			PrevBlockHash: types.RandomHash(),
			Timestamp:     time.Now().UnixNano(),
			Height:        39,
			Nonce:         123,
		},
		Transactions: []Transaction{},
	}
	hash := b.Hash()
	fmt.Println(hash)
	assert.False(t, hash.IsZero())
	// assert.NotNil(t, hash)
}
