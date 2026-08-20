package core

import (
	"fmt"
	"testing"

	"github.com/Wi1low/chainline/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	data := []byte("hello world")
	tx := Transaction{
		Data: data,
	}

	fmt.Printf("sign before: %+v\n", tx)

	privkey := crypto.GeneratePrivateKey()
	assert.Nil(t, tx.Sign(privkey))

	fmt.Printf("sign after: %+v\n", tx)
}
func TestVerifyTransaction(t *testing.T) {
	privkey := crypto.GeneratePrivateKey()

	data := []byte("hello world")
	tx := Transaction{
		Data: data,
	}
	assert.Nil(t, tx.Sign(privkey))
	assert.Nil(t, tx.Verify())

	otherPrikey := crypto.GeneratePrivateKey()
	tx.From = otherPrikey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privkey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: []byte("hello world"),
	}
	assert.Nil(t, tx.Sign(privkey))
	return tx
}
