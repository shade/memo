package main_node

import (
	"fmt"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/user_stats"
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
	_, errors := transaction.ProcessNotifications()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
	_, errors = transaction.UpdateRootTxHashes()
	for _, err := range errors {
		fmt.Println(err.Error())
	}
	if n.BlocksSyncComplete && ! n.UserNode {
		go func() {
			err = user_stats.Populate()
			if err != nil {
				jerr.Get("error populating user stats", err).Print()
			}
		}()
	}
	var prevSaved string
	if ! n.BlocksSyncComplete {
		prevSaved = "prev "
	}
	fmt.Printf("Merkle block height: %5d (%s), hashes: %4d (%ssaved: %4d, memos: %4d)\n",
		dbBlock.Height,
		dbBlock.Timestamp.String(),
		len(transactionHashes),
		prevSaved,
		n.AllTxnsFound,
		n.MemoTxnsFound,
	)
	if dbBlock.Height == n.NodeStatus.HeightChecked+1 {
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
		queueMerkleBlocks(n, false)
	}
}

func queueMerkleBlocks(n *Node, first bool) {
	if n.BlocksQueued != 0 {
		return
	}
	if n.NodeStatus.HeightChecked < MinCheckHeight {
		n.NodeStatus.HeightChecked = MinCheckHeight
	}
	var initialHeight = n.NodeStatus.HeightChecked
	if ! first || n.BlocksSyncComplete {
		initialHeight++
	}
	blocks, err := db.GetBlocksInHeightRange(initialHeight, initialHeight+1999)
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
	if n.BlocksQueued > 1 {
		fmt.Printf("Blocks queued: %d\n", n.BlocksQueued)
	}
}
