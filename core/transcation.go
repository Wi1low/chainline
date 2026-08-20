package core

import (
	"fmt"
	"time"

	"github.com/Wi1low/chainline/crypto"
	"github.com/Wi1low/chainline/types"
)

type Transaction struct {
	Data []byte
	// every transaction has a signature from sender
	From      crypto.PublicKey
	Signature *crypto.Signature

	// cached version of the tx data hash
	hash types.Hash

	// createdAt is the timestamp of when this tx is first seen locally
	createdAt int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:      data,
		createdAt: time.Now().UnixNano(),
	}
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if !tx.hash.IsZero() {
		return tx.hash
	}
	return hasher.Hash(tx)
}
func (tx *Transaction) CreatedAt() int64 {
	return tx.createdAt
}
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("transaction signature is invalid")
	}
	return nil
}
func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}
