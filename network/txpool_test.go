package network

import (
	"strconv"
	"testing"

	"github.com/Wi1low/chainline/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	pool := NewTxPool()

	assert.Equal(t, pool.Len(), 0)
}

func TestTxPoolAdd(t *testing.T) {
	p := NewTxPool()

	tx1 := core.NewTransaction([]byte("tx1"))
	assert.Nil(t, p.Add(tx1))

	assert.Equal(t, p.Len(), 1)

	p.Flush()

	assert.Equal(t, p.Len(), 0)
}

func TestSortTransactions(t *testing.T) {
	pool := NewTxPool()

	txLen := 1000

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		assert.Nil(t, pool.Add(tx))
	}

	assert.Equal(t, txLen, pool.Len())

	txx := pool.Transactions()

	for i := 0; i < len(txx)-1; i++ {
		// create quickly every second
		assert.True(t, txx[i].CreatedAt() <= txx[i+1].CreatedAt())
	}
}
