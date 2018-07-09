package memo

import (
	"github.com/memocash/memo/app/bitcoin/wallet"
)

type Output struct {
	Address wallet.Address
	Amount  int64
	Type    OutputType
	RefData []byte
	Data    []byte
}
