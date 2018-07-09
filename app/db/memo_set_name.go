package db

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"html"
	"time"
)

type MemoSetName struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte `gorm:"size:500"`
	Address    string
	Name       string `gorm:"size:500"`
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoSetName) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoSetName) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoSetName) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoSetName) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoSetName) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoSetNameById(id uint) (*MemoSetName, error) {
	var memoSetName MemoSetName
	err := find(&memoSetName, MemoSetName{
		Id: id,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo set name", err)
	}
	return &memoSetName, nil
}

func GetMemoSetName(txHash []byte) (*MemoSetName, error) {
	var memoSetName MemoSetName
	err := find(&memoSetName, MemoSetName{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo set name", err)
	}
	return &memoSetName, nil
}

func GetSetNamesByTxHashes(txHashes [][]byte) ([]*MemoSetName, error) {
	var memoSetNames []*MemoSetName
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	result := db.
		Where("tx_hash IN (?)", txHashes).
		Find(&memoSetNames)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo set names", result.Error)
	}
	return memoSetNames, nil
}

func GetNameForPkHash(pkHash []byte) (*MemoSetName, error) {
	names, err := GetNamesForPkHashes([][]byte{pkHash})
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hash", err)
	}
	if len(names) == 0 {
		return nil, nil
	}
	return names[0], nil
}

func GetNamesForPkHashes(pkHashes [][]byte) ([]*MemoSetName, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSelect := "JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_set_names" +
		"	GROUP BY pk_hash" +
		") sq ON (sq.id = memo_set_names.id)"
	var memoSetNames []*MemoSetName
	result := db.
		Preload(BlockTable).
		Joins(joinSelect).
		Where("pk_hash IN (?)", pkHashes).
		Find(&memoSetNames)
	if result.Error != nil {
		return nil, jerr.Get("error getting set names", result.Error)
	}
	return memoSetNames, nil
}

func GetUniqueMemoAPkHashesMatchName(searchString string, offset int) ([][]byte, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSelect := "JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_set_names" +
		"	GROUP BY pk_hash" +
		") sq ON (sq.id = memo_set_names.id)"
	rows, err := db.
		Table("memo_set_names").
		Select("DISTINCT(pk_hash)").
		Joins(joinSelect).
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", searchString)).
		Limit(25).
		Offset(offset).
		Rows()
	if err != nil {
		return nil, jerr.Get("error getting distinct pk hashes", err)
	}
	defer rows.Close()
	var pkHashes [][]byte
	for rows.Next() {
		var pkHash []byte
		err := rows.Scan(&pkHash)
		if err != nil {
			return nil, jerr.Get("error scanning row with pkHash", err)
		}
		pkHashes = append(pkHashes, pkHash)
	}
	return pkHashes, nil
}

func GetCountMemoSetName() (uint, error) {
	cnt, err := count(&MemoSetName{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}

func GetSetNames(offset uint) ([]*MemoSetName, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoSetNames []*MemoSetName
	result := db.
		Preload(BlockTable).
		Limit(25).
		Offset(offset).
		Order("id ASC").
		Find(&memoSetNames)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return memoSetNames, nil
}
