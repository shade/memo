package build

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func TopicMessage(topicName string, message string, privateKey *wallet.PrivateKey) (*wire.MsgTx, error) {
	transactions := []transaction.SpendOutput{{
		Type:    transaction.SpendOutputTypeMemoTopicMessage,
		RefData: []byte(topicName),
		Data:    []byte(message),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo tx", err)
	}
	return tx, nil
}
