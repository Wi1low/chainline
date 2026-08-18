package main

import (
	"time"

	"github.com/Wi1low/chainline/network"
)

func main() {
	local := network.NewLocalTransport("LOCAL")
	remote := network.NewLocalTransport("REMOTE")

	local.Connect(remote)
	remote.Connect(local)

	go func() {
		for {
			remote.SendMsg(local.Addr(), []byte("hello from remote"))
			time.Sleep(1 * time.Second)
		}
	}()

	opt := network.ServerOpt{
		Transports: []network.Transport{local},
	}

	s := network.NewServer(opt)
	s.Start()
}
