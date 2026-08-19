package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeypairSignVerifySuccess(t *testing.T) {
	privatekey := GeneratePrivateKey()

	pubKey := privatekey.PublicKey()

	msg := []byte("hello world")
	sig, err := privatekey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(pubKey, msg))

}
func TestGeneratePrivateKey(t *testing.T) {
	privatekey := GeneratePrivateKey()

	pubKey := privatekey.PublicKey()

	otherPriv := GeneratePrivateKey()
	otherPub := otherPriv.PublicKey()

	msg := []byte("hello world")
	sig, err := privatekey.Sign(msg)
	assert.Nil(t, err)
	assert.False(t, sig.Verify(otherPub, msg))
	assert.False(t, sig.Verify(pubKey, []byte("wrong message")))
}
