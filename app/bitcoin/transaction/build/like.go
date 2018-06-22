package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
)

func Like(likeTxBytes []byte, tip int64, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.SpendOutput{{
		Type: memo.SpendOutputTypeMemoLike,
		Data: likeTxBytes,
	}}
	if tip != 0 {
		if tip < memo.DustMinimumOutput {
			return nil, jerr.New("error tip not above dust limit")
		}
		if tip > 1e8 {
			return nil, jerr.New("error trying to tip too much")
		}
		memoPost, err := db.GetMemoPost(likeTxBytes)
		if err != nil {
			return nil, jerr.Get("error getting memo_post", err)
		}
		transactions = append(transactions, memo.SpendOutput{
			Type:    memo.SpendOutputTypeP2PK,
			Address: memoPost.GetAddress(),
			Amount:  tip,
		})
	}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building like tx", err)
	}
	return tx, nil
}
