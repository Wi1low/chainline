package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Hash [32]uint8

func (h Hash) String() string {
	// hex forms
	// return fmt.Sprintf("%x", h[:])
	return hex.EncodeToString(h.ToSlice())
}
func (h Hash) ToSlice() []byte {
	return h[:]
}
func (h Hash) IsZero() bool {
	return h == Hash{}
}
func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}

	return Hash(b)
}

func RandomBytes(size int) []byte {
	token := make([]byte, size)

	// math/rand is deprecated so that we use crypto/rand
	rand.Read(token)
	return token
}

func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}
