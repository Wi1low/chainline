package main

import (
	"bytes"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/Wi1low/chainline/core"
	"github.com/Wi1low/chainline/crypto"
	"github.com/Wi1low/chainline/network"
)

func main() {
	local := network.NewLocalTransport("LOCAL")
	remote := network.NewLocalTransport("REMOTE")

	local.Connect(remote)
	remote.Connect(local)

	go func() {
		for {
			if err := sendTransaction(remote, local.Addr()); err != nil {
				slog.Error("send tx error", slog.Any("error", err))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	opt := network.ServerOpt{
		Transports: []network.Transport{local},
	}

	s := network.NewServer(opt)
	s.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privkey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(100000000)), 10))
	// msg := network.NewMessage()
	tx := core.NewTransaction(data)
	tx.Sign(privkey)

	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMsg(to, msg.Bytes())
	// return tr.SendMsg(tx.)
}
