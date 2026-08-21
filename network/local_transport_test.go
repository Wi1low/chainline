package network

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("a").(*LocalTransport)
	trb := NewLocalTransport("b").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)

	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestSendMsg(t *testing.T) {
	tra := NewLocalTransport("a").(*LocalTransport)
	trb := NewLocalTransport("b").(*LocalTransport)

	tra.Connect(trb)
	trb.Connect(tra)
	msg := []byte("hello world")
	assert.Nil(t, tra.SendMsg(trb.addr, msg))

	msgCh := <-trb.Consume()
	buf := make([]byte, len(msg))

	n, err := msgCh.Payload.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, n, len(msg))

	assert.Equal(t, buf, msg)
	assert.Equal(t, msgCh.From, tra.addr)
}

func TestBroadcast(t *testing.T) {

	tra := NewLocalTransport("a").(*LocalTransport)
	trb := NewLocalTransport("b").(*LocalTransport)
	trc := NewLocalTransport("c").(*LocalTransport)
	trd := NewLocalTransport("d").(*LocalTransport)

	// tra connect b, c, d
	tra.Connect(trb)
	tra.Connect(trc)
	tra.Connect(trd)

	msg := []byte("hello world")
	assert.Nil(t, tra.Broadcast(msg))
	// trb
	rpcb := <-trb.Consume()
	
	b, err := io.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	// trc
	rpcc := <-trc.Consume()

	c, err := io.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, c, msg)

	// trd
	rpcd := <-trd.Consume()

	d, err := io.ReadAll(rpcd.Payload)
	assert.Nil(t, err)
	assert.Equal(t, d, msg)
}
