package main_node

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
)

func onTx(n *Node, msg *wire.MsgTx) {
	if !n.HeaderSyncComplete || !n.BlocksSyncComplete  {
		return
	}
	savedTxn, memoTxn, err := transaction.ConditionallySaveTransaction(msg, nil, n.UserNode)
	if err != nil {
		jerr.Get("error conditionally saving transaction", err).Print()
	}
	if savedTxn {
		if memoTxn {
			fmt.Printf("Saved unconfirmed memo txn: %s\n", msg.TxHash().String())
		} else {
			fmt.Printf("Saved unconfirmed txn: %s\n", msg.TxHash().String())
		}
	}
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
