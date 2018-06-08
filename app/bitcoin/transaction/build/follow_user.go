package build

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func FollowUser(pkHash []byte, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type: transaction.SpendOutputTypeMemoFollow,
		Data: pkHash,
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo tx", err)
	}
	return tx, nil
}

func UnfollowUser(pkHash []byte, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type: transaction.SpendOutputTypeMemoUnfollow,
		Data: pkHash,
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo tx", err)
	}
	return tx, nil
}
