package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

func ProfilePic(url string, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.Output{{
		Type: memo.OutputTypeMemoSetProfilePic,
		Data: []byte(url),
	}}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building profile pic tx", err)
	}
	return tx, nil
}
