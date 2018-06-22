package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
	"sort"
)

func Poll(pollType memo.PollType, question string, options []string, privateKey *wallet.PrivateKey) ([]*memo.Tx, error) {
	var transactionType memo.SpendOutputType
	switch memo.PollType(pollType) {
	case memo.PollTypeOne:
		transactionType = memo.SpendOutputTypeMemoPollQuestionSingle
	case memo.PollTypeAny:
		transactionType = memo.SpendOutputTypeMemoPollQuestionMulti
	default:
		return nil, jerr.New("invalid poll type")
	}

	spendableTxOuts, err := db.GetSpendableTransactionOutputsForPkHash(privateKey.GetPublicKey().GetAddress().GetScriptAddress())
	if err != nil {
		return nil, jerr.Get("error getting spendable tx outs", err)
	}
	sort.Sort(db.TxOutSortByValue(spendableTxOuts))

	var memoTxns []*memo.Tx
	memoTx, spendableTxOuts, err := buildWithTxOuts([]memo.SpendOutput{{
		Type:    transactionType,
		Data:    []byte(question),
		RefData: []byte{byte(len(options))},
	}}, spendableTxOuts, privateKey)
	if err != nil {
		return nil, jerr.Get("error creating tx", err)
	}
	memoTxns = append(memoTxns, memoTx)

	memoTxHash := memoTx.MsgTx.TxHash()
	var questionTxHashBytes = memoTxHash.CloneBytes()
	for _, option := range options {
		memoTx, spendableTxOuts, err = buildWithTxOuts([]memo.SpendOutput{{
			Type:    memo.SpendOutputTypeMemoPollOption,
			Data:    []byte(option),
			RefData: []byte(questionTxHashBytes),
		}}, spendableTxOuts, privateKey)
		if err != nil {
			return nil, jerr.Get("error creating tx", err)
		}
		memoTxns = append(memoTxns, memoTx)
	}

	return memoTxns, nil
}
