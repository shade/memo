package main_node

import (
	"fmt"
	"github.com/jchavannes/bchutil/bloom"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/db"
)

func setBloomFilters(n *Node) {
	if n.UserNode {
		allKeys, err := db.GetAllKeys()
		if err != nil {
			jerr.Get("error getting keys from db", err).Print()
			return
		}
		if len(allKeys) == n.PreviousFilterSize {
			return
		}
		n.PreviousFilterSize = len(allKeys)
		fmt.Printf("Setting bloom filter (keys: %d)...\n", len(allKeys))
		bloomFilter := bloom.NewFilter(uint32(len(allKeys)*2), 0, 0, wire.BloomUpdateNone)
		for _, key := range allKeys {
			bloomFilter.Add(key.GetAddress().GetScriptAddress())
			bloomFilter.Add(key.GetPublicKey().GetSerialized())
		}
		n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	} else {
		codes := memo.GetAllCodes()
		if len(codes) == n.PreviousFilterSize {
			return
		}
		n.PreviousFilterSize = len(codes)
		fmt.Printf("Setting bloom filter (codes: %d)...\n", len(codes))
		bloomFilter := bloom.NewFilter(uint32(len(codes)), 0, 0, wire.BloomUpdateNone)
		for _, code := range codes {
			bloomFilter.Add(code)
		}
		n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	}
}
