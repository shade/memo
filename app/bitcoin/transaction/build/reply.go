package build

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func MemoReply(txHashBytes []byte, message string, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type:    transaction.SpendOutputTypeMemoReply,
		RefData: txHashBytes,
		Data:    []byte(message),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo reply tx", err)
	}
	return tx, nil
}
