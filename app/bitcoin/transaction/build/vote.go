package build

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
)

func Vote(pollTxBytes []byte, message string, tip int64, privateKey *wallet.PrivateKey) (*memo.Tx, error) {
	transactions := []memo.Output{{
		Type:    memo.OutputTypeMemoPollVote,
		Data:    pollTxBytes,
		RefData: []byte(message),
	}}
	if tip != 0 {
		if tip < memo.DustMinimumOutput {
			return nil, jerr.New("error tip not above dust limit")
		}
		if tip > 1e8 {
			return nil, jerr.New("error trying to tip too much")
		}
		memoPollOption, err := db.GetMemoPollOption(pollTxBytes)
		if err != nil {
			return nil, jerr.Get("error getting memo poll option", err)
		}
		transactions = append(transactions, memo.Output{
			Type:    memo.OutputTypeP2PK,
			Address: memoPollOption.GetAddress(),
			Amount:  tip,
		})
	}
	tx, err := Build(transactions, privateKey)
	if err != nil {
		return nil, jerr.Get("error building vote tx", err)
	}
	return tx, nil
}
