package main_node

import (
	"fmt"
	"github.com/jchavannes/bchutil/bloom"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/db"
)

var previousFilterSize int

func setBloomFilters(n *Node) {
	allKeys, err := db.GetAllKeys()
	if err != nil {
		jerr.Get("error getting keys from db", err).Print()
		return
	}
	if len(allKeys) == previousFilterSize {
		return
	}
	previousFilterSize = len(allKeys)
	codes := memo.GetAllCodes()
	fmt.Printf("Setting bloom filter (keys: %d, codes: %d)...\n", len(allKeys), len(codes))
	bloomFilter := bloom.NewFilter(uint32(len(allKeys)*2), 0, 0, wire.BloomUpdateNone)
	for _, key := range allKeys {
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
		bloomFilter.Add(key.GetPublicKey().GetSerialized())
	}
	for _, code := range codes {
		bloomFilter.Add(code)
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
}
