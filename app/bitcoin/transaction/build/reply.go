package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func MemoReply(txHashBytes []byte, message string, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.SpendOutput{{
		Type:    memo.SpendOutputTypeMemoReply,
		RefData: txHashBytes,
		Data:    []byte(message),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo reply tx", err)
	}
	return tx, nil
}
