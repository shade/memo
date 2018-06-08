package build

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func FollowTopic(topicName string, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type: transaction.SpendOutputTypeMemoTopicFollow,
		Data: []byte(topicName),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building topic follow tx", err)
	}
	return tx, nil
}

func UnfollowTopic(topicName string, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type: transaction.SpendOutputTypeMemoTopicUnfollow,
		Data: []byte(topicName),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building topic unfollow tx", err)
	}
	return tx, nil
}
