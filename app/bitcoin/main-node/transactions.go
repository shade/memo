package main_node

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
)

func onTx(n *Node, msg *wire.MsgTx) {
	block := findHashBlock([]map[string]*db.Block{n.BlockHashes, n.PrevBlockHashes}, msg.TxHash())
	if (!n.HeaderSyncComplete || !n.BlocksSyncComplete) && block == nil {
		return
	}
	savedTxn, memoTxn, err := transaction.ConditionallySaveTransaction(msg, block)
	if err != nil {
		jerr.Get("error conditionally saving transaction", err).Print()
	}
	if savedTxn {
		n.AllTxnsFound++
		if memoTxn {
			n.MemoTxnsFound++
			if block == nil {
				fmt.Printf("Saved unconfirmed memo txn: %s\n", msg.TxHash().String())
			}
		} else if block == nil {
			fmt.Printf("Saved unconfirmed txn: %s\n", msg.TxHash().String())
		}
	}
}

func findHashBlock(blockHashes []map[string]*db.Block, hash chainhash.Hash) *db.Block {
	for _, hashMap := range blockHashes {
		for hashString, block := range hashMap {
			if hashString == hash.String() {
				return block
			}
		}
	}
	return nil
}

func getTransaction(n *Node, txId chainhash.Hash) {
	msgGetData := wire.NewMsgGetData()
	err := msgGetData.AddInvVect(&wire.InvVect{
		Type: wire.InvTypeTx,
		Hash: txId,
	})
	if err != nil {
		jerr.Get("error adding invVect: %s\n", err).Print()
		return
	}
	n.Peer.QueueMessage(msgGetData, nil)
}
