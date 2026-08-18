package network

import (
	"fmt"
	"time"
)

type ServerOpt struct {
	Transports []Transport
}

type Server struct {
	ServerOpt
	rpcCh chan RPC
	quit  chan struct{}
}

func NewServer(opt ServerOpt) *Server {
	return &Server{
		ServerOpt: opt,
		rpcCh:     make(chan RPC),
		quit:      make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quit:
			break free
		case <-ticker.C:
			fmt.Println("tick tick tick...")
		}
	}
	fmt.Println("Server Shutdown")
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
