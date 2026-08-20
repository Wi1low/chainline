package network

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Wi1low/chainline/core"
	"github.com/Wi1low/chainline/crypto"
)

var defaultBlockTime = 5 * time.Second

type ServerOpt struct {
	Transports []Transport
	PrivateKey *crypto.PrivateKey
	BlockTime  time.Duration
}

type Server struct {
	ServerOpt
	memPool *TxPool
	// isValidator indicates whether the server is a validator
	blockTime   time.Duration
	isValidator bool
	rpcCh       chan RPC
	quit        chan struct{}
}

func NewServer(opt ServerOpt) *Server {
	if opt.BlockTime == time.Duration(0) {
		opt.BlockTime = defaultBlockTime
	}
	return &Server{
		ServerOpt:   opt,
		memPool:     NewTxPool(),
		blockTime:   opt.BlockTime,
		isValidator: opt.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quit:        make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quit:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}
	fmt.Println("Server Shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {

	if err := tx.Verify(); err != nil {
		return err
	}
	hash := tx.Hash(core.TxHasher)
	if s.memPool.Has(hash) {
		slog.Info("transaction already in pool",
			slog.Any("tx hash", hash))
		return nil
	}

	slog.Info("adding new tx to the mempool",
		slog.Any("tx hash", hash))

	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("create a new block")

	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
func (s *Server) Stop() {
	s.quit <- struct{}{}
}
