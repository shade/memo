package main_node

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/bchutil"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/user_stats"
	"time"
	"github.com/memocash/memo/app/metric"
)

const MinCheckHeight = 525000

func onBlock(n *Node, msg *wire.MsgBlock) {
	block := bchutil.NewBlock(msg)
	dbBlock, err := db.GetBlockByHash(*block.Hash())
	if err != nil {
		jerr.Getf(err, "error getting dbBlock (%s)", block.Hash().String()).Print()
		return
	}
	var memosSaved int
	var txnsSaved int
	for _, txn := range block.Transactions() {
		saveStart := time.Now()
		savedTxn, savedMemo, err := transaction.ConditionallySaveTransaction(txn.MsgTx(), dbBlock)
		if err != nil {
			jerr.Getf(err, "error conditionally saving transaction: %s", txn.Hash().String()).Print()
			continue
		}
		if savedTxn {
			txnsSaved++
		}
		if savedMemo {
			memosSaved++
		}
		metric.AddTransactionSaveTime(time.Since(saveStart))
	}
	_, errors := transaction.ProcessNotifications()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
	_, errors = transaction.UpdateRootTxHashes()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
	go func() {
		if n.BlocksSyncComplete {
			err = user_stats.Populate()
			if err != nil {
				jerr.Get("error populating user stats", err).Print()
			}
		}
	}()
	fmt.Printf("Block - height: %5d (%s), found: %4d, saved: %4d, memos: %4d\n",
		dbBlock.Height,
		dbBlock.Timestamp.String(),
		len(block.Transactions()),
		txnsSaved,
		memosSaved,
	)
	if dbBlock.Height == n.NodeStatus.HeightChecked + 1 {
		n.NodeStatus.HeightChecked = dbBlock.Height
		err = n.NodeStatus.Save()
	}
	if err != nil {
		jerr.Get("error saving node status", err).Print()
		return
	}
	n.BlocksQueued--
	if n.BlocksQueued == 0 {
		queueBlocks(n)
	}
}

func queueBlocks(n *Node) {
	if n.BlocksQueued != 0 {
		return
	}
	if n.NodeStatus.HeightChecked < MinCheckHeight {
		n.NodeStatus.HeightChecked = MinCheckHeight
	}
	blocks, err := db.GetBlocksInHeightRange(n.NodeStatus.HeightChecked+1, n.NodeStatus.HeightChecked+2000)
	if err != nil {
		jerr.Get("error getting blocks in height range", err).Print()
		return
	}
	if len(blocks) == 0 {
		if ! n.BlocksSyncComplete {
			n.BlocksSyncComplete = true
			transaction.DisableBatchPostProcessing()
			fmt.Println("Block sync complete")
			queueMempool(n)
		}
		return
	}
	msgGetData := wire.NewMsgGetData()
	for _, block := range blocks {
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			jerr.Get("error adding inventory vector: %s\n", err).Print()
			return
		}
	}
	n.Peer.QueueMessage(msgGetData, nil)
	n.BlocksQueued += len(msgGetData.InvList)
	if n.BlocksQueued > 1 {
		fmt.Printf("Blocks queued: %d\n", n.BlocksQueued)
	}
}

func getBlock(n *Node, hash chainhash.Hash) {
	getBlocks := wire.NewMsgGetBlocks(&hash)
	n.Peer.QueueMessage(getBlocks, nil)
}
