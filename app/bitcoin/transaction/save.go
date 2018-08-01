package transaction

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"time"
	"github.com/memocash/memo/app/metric"
)

func ConditionallySaveTransaction(msg *wire.MsgTx, dbBlock *db.Block, userNode bool) (bool, bool, error) {
	saveStart := time.Now()
	dbTxn, err := db.ConvertMsgToTransaction(msg)
	if err != nil {
		// Don't log, lots of mal-formed txns
		// jerr.Get("error converting msg to db transaction", err)
		return false, false, nil
	}
	memoOutput, err := GetMemoOutputIfExists(dbTxn)
	if err != nil {
		return false, false, jerr.Get("error getting memo output", err)
	}
	pkHashes := GetPkHashesFromTxn(dbTxn)
	var savingMemo bool
	if (memoOutput == nil && !userNode) || (memoOutput != nil && userNode) {
		return false, false, nil
	} else if memoOutput == nil {
		watched, err := db.ContainsWatchedPkHash(pkHashes)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return false, false, jerr.Get("error checking db for watched addresses", err)
		}
		if ! watched {
			return false, false, nil
		}
	} else {
		savingMemo = true
	}
	err = SaveTransaction(dbTxn, dbBlock)
	if err != nil {
		return false, false, jerr.Get("error saving transaction", err)
	}
	err = ClearCaches(pkHashes)
	if err != nil {
		return false, false, jerr.Get("error clearing transaction caches", err)
	}
	err = metric.AddTransactionSaveTime(time.Since(saveStart))
	if err != nil {
		return false, false, jerr.Get("error add transaction save time metric", err)
	}
	return true, savingMemo, nil
}

func SaveTransaction(txn *db.Transaction, block *db.Block) error {
	hash := txn.GetChainHash()
	existingTxn, err := db.GetTransactionByHash(hash.CloneBytes())
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting transaction from db", err)
	}
	if existingTxn != nil {
		if block != nil && existingTxn.BlockId == 0 {
			err = updateTxn(existingTxn, block)
			if err != nil {
				return jerr.Get("error updating transaction", err)
			}
		}
	} else {
		err = newTxn(txn, block)
		if err != nil {
			return jerr.Get("error saving new transaction", err)
		}
	}
	memoOutput, err := GetMemoOutputIfExists(txn)
	if err != nil {
		return jerr.Get("error getting memo output", err)
	}
	if memoOutput != nil {
		err = SaveMemo(txn, memoOutput, block)
		if err != nil {
			return jerr.Get("error saving memos", err)
		}
	}
	return nil
}

func ClearCaches(pkHashes [][]byte) error {
	for _, pkHash := range pkHashes {
		err := cache.ClearBalance(pkHash)
		if err != nil && ! cache.IsMissError(err) {
			return jerr.Get("error clearing balance cache", err)
		}
	}
	return nil
}

func updateTxn(txn *db.Transaction, block *db.Block) error {
	if block != nil {
		txn.BlockId = block.Id
	}
	err := txn.Save()
	if err != nil {
		return jerr.Get("error saving updated transaction", err)
	}
	return nil
}

func newTxn(txn *db.Transaction, block *db.Block) error {
	var blockId uint
	if block != nil {
		blockId = block.Id
	}
	txn.BlockId = blockId
	for _, in := range txn.TxIn {
		err := updateExistingOutput(in)
		if err != nil {
			return jerr.Get("error updating existing outputs", err)
		}
	}

	for _, out := range txn.TxOut {
		err := updateExistingInput(out, txn.Hash)
		if err != nil {
			return jerr.Get("error updating existing inputs", err)
		}
	}

	err := txn.Save()
	if err != nil {
		return jerr.Get("error saving txn", err)
	}

	return nil
}

func updateExistingOutput(in *db.TransactionIn) error {
	existingTxOut, err := db.GetTransactionOutput(in.PreviousOutPointHash, in.PreviousOutPointIndex)
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			return nil
		}
		return jerr.Get("error getting existing output", err)
	}
	existingTxOut.TxnInHashString = in.HashString
	err = existingTxOut.Save()
	if err != nil {
		return jerr.Get("error saving existing transaction output", err)
	}
	return nil
}

func updateExistingInput(out *db.TransactionOut, txHash []byte) error {
	existingInput, err := db.GetTransactionInput(txHash, out.Index)
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			return nil
		}
		return jerr.Get("error getting existing input", err)
	}
	out.TxnInHashString = existingInput.HashString
	return nil
}
