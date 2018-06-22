package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func SetProfileText(profileText string, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.SpendOutput{{
		Type: memo.SpendOutputTypeMemoSetProfile,
		Data: []byte(profileText),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building memo set profile text tx", err)
	}
	return tx, nil
}
