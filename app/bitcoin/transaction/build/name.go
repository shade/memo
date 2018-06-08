package build

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func SetName(name string, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type: transaction.SpendOutputTypeMemoSetName,
		Data: []byte(name),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo set name tx", err)
	}
	return tx, nil
}
