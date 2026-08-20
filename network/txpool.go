package network

import (
	"sort"

	"github.com/Wi1low/chainline/core"
	"github.com/Wi1low/chainline/types"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))

	i := 0

	for _, tx := range txMap {
		txx[i] = tx
		i++
	}

	s := &TxMapSorter{transactions: txx}

	sort.Sort(s)

	return s
}
func (tms *TxMapSorter) Len() int {
	return len(tms.transactions)
}
func (tms *TxMapSorter) Less(i, j int) bool {
	return tms.transactions[i].CreatedAt() < tms.transactions[j].CreatedAt()
}

func (tms *TxMapSorter) Swap(i, j int) {
	tms.transactions[i], tms.transactions[j] = tms.transactions[j], tms.transactions[i]
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (tp *TxPool) Transactions() []*core.Transaction {
	return NewTxMapSorter(tp.transactions).transactions
}

// Add adds an transaction to the pool, the caller is repsonseible checking if the
// tx already exist
func (tp *TxPool) Add(tx *core.Transaction) error {
	thash := tx.Hash(core.TxHasher)

	// if tp.Has(thash) {
	// 	return nil
	// }

	tp.transactions[thash] = tx
	return nil
}
func (tp *TxPool) Has(txhash types.Hash) bool {
	_, ok := tp.transactions[txhash]
	return ok
}
func (tp *TxPool) Len() int {
	return len(tp.transactions)
}

func (tp *TxPool) Flush() {
	tp.transactions = make(map[types.Hash]*core.Transaction)
}
