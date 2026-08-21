package network

type NetAddr string

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	Broadcast([]byte) error
	SendMsg(NetAddr, []byte) error
	Addr() NetAddr
}
