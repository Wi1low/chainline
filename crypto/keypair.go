package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/Wi1low/chainline/types"
)

type PrivateKey struct {
	Key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
}

func GeneratePrivateKey() PrivateKey {
	// TODO: implement
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{Key: key}
}
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: &k.Key.PublicKey}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) ToSlice() []byte {
	// ecdh.(k.key.Curve, k.key.X, k.key.Y)
	return elliptic.MarshalCompressed(k.Key.Curve, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	// 256 bits (32 bytes)
	h := sha256.Sum256(k.ToSlice())
	// 32 - 20 = 12
	// h[12:] is the last 20 bytes of the hash
	return types.AddressFromBytes(h[len(h)-20:])
}

// GobEncode serializes the public key as a compressed point so that gob does
// not need to encode the elliptic.Curve field (which has no exported fields).
func (k PublicKey) GobEncode() ([]byte, error) {
	return elliptic.MarshalCompressed(k.Key.Curve, k.Key.X, k.Key.Y), nil
}

// GobDecode reconstructs the public key from a compressed point.
func (k *PublicKey) GobDecode(data []byte) error {
	curve := elliptic.P256()
	x, y := elliptic.UnmarshalCompressed(curve, data)
	if x == nil {
		return fmt.Errorf("invalid public key data")
	}
	k.Key = &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return nil
}

type Signature struct {
	R, S *big.Int
}

func (sig Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S)
}
