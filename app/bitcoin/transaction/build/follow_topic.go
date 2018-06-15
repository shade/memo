package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func FollowTopic(topicName string, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.SpendOutput{{
		Type: memo.SpendOutputTypeMemoTopicFollow,
		Data: []byte(topicName),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building topic follow tx", err)
	}
	return tx, nil
}

func UnfollowTopic(topicName string, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.SpendOutput{{
		Type: memo.SpendOutputTypeMemoTopicUnfollow,
		Data: []byte(topicName),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building topic unfollow tx", err)
	}
	return tx, nil
}
