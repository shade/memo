package db

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db/obj"
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

func GetMemoCohortStats() ([]obj.MemoCohortStat, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_tests").
		Select("COUNT(*) AS num_posts, " +
		"COUNT(DISTINCT memo_tests.pk_hash) AS num_users," +
		"DATE(DATE_FORMAT(user_stats.first_post, '%Y-%m-01')) AS cohort," +
		"DATE(DATE_FORMAT(`timestamp`, '%Y-%m-%d')) AS date").
		Joins("JOIN blocks ON (memo_tests.block_id = blocks.id)").
		Joins("JOIN user_stats ON (memo_tests.pk_hash = user_stats.pk_hash)").
		Where("HEX(memo_tests.pk_script) REGEXP '6A026D(0|1).*'").
		Group("date, cohort").
		Order("date ASC")
	var memoCohortStats []obj.MemoCohortStat
	result := query.Find(&memoCohortStats)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo stats", result.Error)
	}
	return memoCohortStats, nil
}

func GetUserStats() ([]obj.UserStat, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSql := "" +
		"LEFT JOIN (" +
		"	SELECT" +
		"		follow_pk_hash," +
		"		COALESCE(SUM(IF(unfollow=0, 1, 0)), 0) AS num_followers" +
		"	FROM memo_follows" +
		"	JOIN (" +
		"		SELECT MAX(id) AS id" +
		"		FROM memo_follows" +
		"		GROUP BY pk_hash, follow_pk_hash" +
		"	) sq ON (sq.id = memo_follows.id)" +
		"	GROUP BY follow_pk_hash" +
		") follows ON (memo_tests.pk_hash = follows.follow_pk_hash)"
	query := db.
		Table("memo_tests").
		Select("memo_tests.pk_hash, " +
		"COUNT(*) AS num_posts," +
		"follows.num_followers AS num_followers," +
		"MIN(`timestamp`) AS first_post," +
		"MAX(`timestamp`) AS last_post").
		Joins("JOIN blocks ON (memo_tests.block_id = blocks.id)").
		Joins(joinSql).
		Where("HEX(memo_tests.pk_script) REGEXP '6A026D(0|1).*'").
		Group("memo_tests.pk_hash")
	var userStats []obj.UserStat
	result := query.Find(&userStats)
	if result.Error != nil {
		return nil, jerr.Get("error getting user stats", result.Error)
	}
	return userStats, nil
}
