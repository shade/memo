package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func FollowUser(pkHash []byte, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.Output{{
		Type: memo.OutputTypeMemoFollow,
		Data: pkHash,
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo tx", err)
	}
	return tx, nil
}

func UnfollowUser(pkHash []byte, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.Output{{
		Type: memo.OutputTypeMemoUnfollow,
		Data: pkHash,
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo tx", err)
	}
	return tx, nil
}
