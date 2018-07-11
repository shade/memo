package db

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"html"
	"time"
)

type MemoTest struct {
	Id        uint   `gorm:"primary_key"`
	TxHash    []byte `gorm:"unique;size:50"`
	PkHash    []byte
	PkScript  []byte `gorm:"size:500"`
	Address   string
	BlockId   uint
	Block     *Block
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m MemoTest) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoTest) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo test", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoTest) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo test", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoTest) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoTest) GetCode() string {
	if len(m.PkScript) < 5 {
		return ""
	}
	return hex.EncodeToString(m.PkScript[2:4])
}

func GetMemoTest(txHash []byte) (*MemoTest, error) {
	var memoTest MemoTest
	err := find(&memoTest, MemoTest{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo test", err)
	}
	return &memoTest, nil
}

type MemoStat struct {
	Date     time.Time
	NumPosts int
	NumUsers int
}

func GetStats() ([]MemoStat, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_tests").
		Select("COUNT(*) AS num_posts, " +
		"COUNT(DISTINCT pk_hash) AS num_users," +
		"DATE(DATE_FORMAT(`timestamp`, '%Y-%m-%d')) AS date").
		Joins("JOIN blocks ON (memo_tests.block_id = blocks.id)").
		Group("date").
		Order("date ASC")
	var memoStats []MemoStat
	result := query.Find(&memoStats)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo stats", result.Error)
	}
	return memoStats, nil
}
