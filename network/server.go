package network

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"

	"github.com/Wi1low/chainline/core"
	"github.com/Wi1low/chainline/crypto"
)

var defaultBlockTime = 5 * time.Second

type ServerOpt struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	PrivateKey    *crypto.PrivateKey
	BlockTime     time.Duration
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

	if opt.RPCDecodeFunc == nil {
		opt.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		ServerOpt:   opt,
		memPool:     NewTxPool(),
		blockTime:   opt.BlockTime,
		isValidator: opt.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quit:        make(chan struct{}, 1),
	}

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(s.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)

			if err != nil {
				slog.Error("rpc error", slog.Any("error", err))
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				slog.Error("rpc processor transaction error", slog.Any("error", err))
			}

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

func (s *Server) broadcast(msg []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(msg); err != nil {
			return err
		}

	}
	return nil
}
func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())

}
func (s *Server) ProcessMessage(dmsg *DecodedMessage) error {

	switch t := dmsg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(dmsg.From, t)
	}
	return nil
}

func (s *Server) processTransaction(from NetAddr, tx *core.Transaction) error {

	hash := tx.Hash(core.TxHasher)
	if s.memPool.Has(hash) {
		slog.Info("transaction already in pool",
			slog.Any("tx hash", hash))
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	slog.Info("adding new tx to the mempool",
		slog.Any("tx hash", hash),
		slog.Int("mempool length", s.memPool.Len()))

	go s.broadcastTx(tx)

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
