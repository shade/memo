package main_node

import (
	"fmt"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"time"
)

func onVerAck(n *Node, msg *wire.MsgVerAck) {
	setBloomFilters(n)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			setBloomFilters(n)
		}
	}()
	block, err := db.GetRecentBlock()
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	sendGetHeaders(n, block.GetChainhash())
}
