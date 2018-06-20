package main_node

import (
	"fmt"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
)

func onMerkleBlock(n *Node, msg *wire.MsgMerkleBlock) {
	dbBlock, err := db.GetBlockByHash(msg.Header.BlockHash())
	if err != nil {
		jerr.Getf(err, "error getting dbBlock (%s)", msg.Header.BlockHash().String()).Print()
		return
	}

	transactionHashes := transaction.GetTransactionsFromMerkleBlock(msg)
	for _, transactionHash := range transactionHashes {
		n.BlockHashes[transactionHash.GetTxId().String()] = dbBlock
	}
	fmt.Printf("Got merkle block - height: %5d, timestamp: %s, hashes: %3d (Prev block - saved: %5d, memos: %5d)\n",
		dbBlock.Height,
		dbBlock.Timestamp.String(),
		len(transactionHashes),
		n.AllTxnsFound,
		n.MemoTxnsFound,
	)
	if dbBlock.Height == n.NodeStatus.HeightChecked + 1 {
		n.NodeStatus.HeightChecked = dbBlock.Height
		err = n.NodeStatus.Save()
	}
	if err != nil {
		jerr.Get("error saving node status", err).Print()
		return
	}
	n.AllTxnsFound = 0
	n.MemoTxnsFound = 0

	n.BlocksQueued--
	if n.BlocksQueued == 0 {
		queueMerkleBlocks(n)
	}
}

func queueMerkleBlocks(n *Node) {
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
			fmt.Println("Block sync complete")
			queueMempool(n)
		}
		return
	}
	msgGetData := wire.NewMsgGetData()
	for _, block := range blocks {
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeFilteredBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			jerr.Get("error adding inventory vector: %s\n", err).Print()
			return
		}
	}
	n.Peer.QueueMessage(msgGetData, nil)
	n.PrevBlockHashes = n.BlockHashes
	n.BlockHashes = make(map[string]*db.Block)
	n.BlocksQueued += len(msgGetData.InvList)
	fmt.Printf("Blocks queued: %d\n", n.BlocksQueued)
}
