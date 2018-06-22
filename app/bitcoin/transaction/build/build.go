package build

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
	"sort"
)

func Build(spendOutputs []memo.SpendOutput, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	spendableTxOuts, err := db.GetSpendableTransactionOutputsForPkHash(privateKey.GetPublicKey().GetAddress().GetScriptAddress())
	if err != nil {
		return nil, jerr.Get("error getting spendable tx outs", err)
	}
	sort.Sort(db.TxOutSortByValue(spendableTxOuts))
	memoTx, _, err := buildWithTxOuts(spendOutputs, spendableTxOuts, privateKey)
	if err != nil {
		return nil, jerr.Get("error creating tx", err)
	}
	return memoTx, nil
}

func buildWithTxOuts(spendOutputs []memo.SpendOutput, spendableTxOuts []*db.TransactionOut, privateKey *wallet.PrivateKey) (*memo.Tx, []*db.TransactionOut, error) {
	var minInput = int64(memo.BaseTxFee + memo.InputFeeP2PKH + memo.OutputFeeP2PKH + memo.DustMinimumOutput)

	for _, spendOutput := range spendOutputs {
		switch spendOutput.Type {
		case memo.SpendOutputTypeP2PK:
			minInput += memo.OutputFeeP2PKH + spendOutput.Amount
		default:
			outputFee, err := getMemoOutputFee(spendOutput)
			if err != nil {
				return nil, nil, jerr.Get("error getting memo output fee", err)
			}
			minInput += outputFee
		}
	}

	var txOutsToUse []*db.TransactionOut
	var totalInputValue int64
	for {
		if len(spendableTxOuts) == 0 {
			return nil, nil, jerr.New("unable to find enough value to spend")
		}
		spendableTxOut := spendableTxOuts[0]
		spendableTxOuts = spendableTxOuts[1:]
		txOutsToUse = append(txOutsToUse, spendableTxOut)
		totalInputValue += spendableTxOut.Value
		if totalInputValue > minInput {
			break
		}
		minInput += memo.InputFeeP2PKH
	}

	var fee = int64(memo.BaseTxFee+len(txOutsToUse)*memo.InputFeeP2PKH) + memo.OutputFeeP2PKH

	var totalOutputValue int64
	for _, spendOutput := range spendOutputs {
		totalOutputValue += spendOutput.Amount
		switch spendOutput.Type {
		case memo.SpendOutputTypeP2PK:
			fee += memo.OutputFeeP2PKH
		default:
			outputFee, err := getMemoOutputFee(spendOutput)
			if err != nil {
				return nil, nil, jerr.Get("error getting memo output fee", err)
			}
			fee += outputFee
		}
	}

	var change = totalInputValue - fee - totalOutputValue
	if change < memo.DustMinimumOutput {
		return nil, nil, jerr.New("change value below dust minimum input")
	}
	address := privateKey.GetPublicKey().GetAddress()
	spendOutputs = append([]memo.SpendOutput{{
		Type:    memo.SpendOutputTypeP2PK,
		Address: address,
		Amount:  change,
	}}, spendOutputs...)

	var tx *wire.MsgTx
	tx, err := transaction.Create(txOutsToUse, privateKey, spendOutputs)
	if err != nil {
		return nil, nil, jerr.Get("error creating tx", err)
	}
	var inputs []*memo.TxInput
	for _, txOut := range txOutsToUse {
		inputs = append(inputs, &memo.TxInput{
			PkHash:      txOut.KeyPkHash,
			Value:       txOut.Value,
			PrevOutHash: txOut.GetHashString(),
		})
	}
	txHash := tx.TxHash()
	var index uint32 = 0
	spendableTxOuts = append([]*db.TransactionOut{{
		TransactionHash: txHash.CloneBytes(),
		PkScript:        tx.TxOut[index].PkScript,
		Index:           index,
		Value:           change,
	}}, spendableTxOuts...)
	return &memo.Tx{
		SelfPkHash: address.GetScriptAddress(),
		MsgTx:      tx,
		Inputs:     inputs,
	}, spendableTxOuts, nil
}

func getMemoOutputFee(spendOutput memo.SpendOutput) (int64, error) {
	switch spendOutput.Type {
	case memo.SpendOutputTypeMemoMessage,
		memo.SpendOutputTypeMemoLike,
		memo.SpendOutputTypeMemoSetName,
		memo.SpendOutputTypeMemoSetProfile,
		memo.SpendOutputTypeMemoSetProfilePic,
		memo.SpendOutputTypeMemoFollow, memo.SpendOutputTypeMemoUnfollow,
		memo.SpendOutputTypeMemoTopicFollow, memo.SpendOutputTypeMemoTopicUnfollow:
		return int64(memo.OutputFeeOpReturn + len(spendOutput.Data)), nil
	case memo.SpendOutputTypeMemoReply,
		memo.SpendOutputTypeMemoTopicMessage,
		memo.SpendOutputTypeMemoPollOption,
		memo.SpendOutputTypeMemoPollVote:
		return int64(memo.OutputFeeOpReturn + len(spendOutput.Data) + memo.OutputOpDataFee + len(spendOutput.RefData)), nil
	case memo.SpendOutputTypeMemoPollQuestionSingle,
		memo.SpendOutputTypeMemoPollQuestionMulti:
		return int64(memo.OutputFeeOpReturn + (memo.OutputOpDataFee+1)*2 + len(spendOutput.Data)), nil
	}
	return 0, jerr.New("unable to get fee for output type")
}
