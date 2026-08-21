package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log/slog"

	"github.com/Wi1low/chainline/core"
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}

type MessageType byte

const (
	MessageTypeTx    MessageType = 0x01
	MessageTypeBlock MessageType = 0x02
)

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(mType MessageType, data []byte) *Message {
	return &Message{Header: mType, Data: data}
}

func (m *Message) Bytes() []byte {
	buf := &bytes.Buffer{}

	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type RPCProcessor interface {
	ProcessMessage(*DecodedMessage) error
}

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err.Error())
	}

	slog.Info("new income msg",
		slog.Any("from", rpc.From),
		slog.Any("type", msg.Header))

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{From: rpc.From, Data: tx}, nil
	// case MessageTypeBlock:
	// 	return d.p.ProcessBlock(rpc)
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}
	// return nil
}
