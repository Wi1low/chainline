package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport, 16),
	}
}

func (lt *LocalTransport) Consume() <-chan RPC {
	return lt.consumeCh
}
func (lt *LocalTransport) Connect(peer Transport) error {
	lt.lock.Lock()
	defer lt.lock.Unlock()

	lt.peers[peer.Addr()] = peer.(*LocalTransport)
	return nil
}
func (lt *LocalTransport) SendMsg(to NetAddr, msg []byte) error {
	lt.lock.RLock()
	defer lt.lock.RUnlock()

	peer, ok := lt.peers[to]
	if !ok {
		return fmt.Errorf("%s could not send message to %s", lt.addr, to)
	}

	peer.consumeCh <- RPC{
		From:    lt.addr,
		Payload: msg,
	}
	return nil
}
func (lt *LocalTransport) Addr() NetAddr {
	return lt.addr
}
