package network

import (
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
	assert.Equal(t, msgCh.Payload, msg)
	assert.Equal(t, msgCh.From, tra.addr)
}
