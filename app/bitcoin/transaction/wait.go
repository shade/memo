package transaction

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/queuer"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/metric"
	"time"
)

const waitTime = 200 * time.Millisecond

func QueueTx(tx *memo.Tx) {
	go func() {
		err := metric.AddMemoBroadcast(tx.Type)
		if err != nil {
			jerr.Get("error adding memo broadcast metric", err).Print()
		}
	}()
	doneChan := make(chan struct{}, 1)
	queuer.Node.Peer.QueueMessage(tx.MsgTx, doneChan)
	<-doneChan
}

func WaitForTx(txHash *chainhash.Hash) error {
	// wait up to 30 seconds
	for i := 0; i < 150; i++ {
		_, err := db.GetTransactionByHash(txHash.CloneBytes())
		if err == nil {
			return nil
		}
		if ! db.IsRecordNotFoundError(err) {
			return jerr.Get("error looking for transaction", err)
		}
		time.Sleep(waitTime)
	}
	return jerr.New("unable to find transaction")
}

func WaitForPic(txHash *chainhash.Hash) error {
	// wait up to 30 seconds
	for i := 0; i < 150; i++ {
		_, err := db.GetMemoSetPic(txHash.CloneBytes())
		if err == nil {
			return nil
		}
		if ! db.IsRecordNotFoundError(err) {
			return jerr.Get("error looking for pic", err)
		}
		time.Sleep(waitTime)
	}
	return jerr.New("unable to find pic")
}
