package db

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"strings"
	"time"
)

var transactionInColumns = []string{
	KeyTable,
	TransactionTable,
	TransactionBlockTbl,
}

type TransactionIn struct {
	Id                    uint            `gorm:"primary_key"`
	Index                 uint
	HashString            string          `gorm:"index:hash_string"`
	TransactionHash       []byte
	Transaction           *Transaction    `gorm:"foreignkey:TransactionHash"`
	KeyPkHash             []byte          `gorm:"index:pk_hash"`
	Key                   *Key            `gorm:"foreignkey:KeyPkHash"`
	PreviousOutPointHash  []byte          `gorm:"unique_index:previous_out"`
	PreviousOutPointIndex uint32          `gorm:"unique_index:previous_out"`
	SignatureScript       []byte
	UnlockString          string
	Sequence              uint32
	TxnOutHashString      string          `gorm:"size:4096"`
	TxnOut                *TransactionOut `gorm:"foreignkey:TxnOutHashString"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (t TransactionIn) GetHashString() string {
	hash, err := chainhash.NewHash(t.TransactionHash)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("in:%s:%d", hash.String(), t.Index)
}

func (t TransactionIn) Save() error {
	result := save(&t)
	if result.Error != nil {
		return jerr.Get("error saving transaction input", result.Error)
	}
	return nil
}

func (t TransactionIn) GetOutPoint() *wire.OutPoint {
	hash, _ := chainhash.NewHash(t.PreviousOutPointHash)
	return wire.NewOutPoint(hash, t.PreviousOutPointIndex)
}

func (t TransactionIn) GetPrevOutPointHash() []byte {
	return t.GetOutPoint().Hash.CloneBytes()
}

func (t TransactionIn) GetPrevOutPointString() string {
	return t.GetOutPoint().String()
}

func (t TransactionIn) HasOut() bool {
	return t.TxnOut != nil
}

func (t TransactionIn) GetPublicKey() []byte {
	split := strings.Split(t.UnlockString, " ")
	if len(split) != 2 {
		return []byte{}
	}
	pubKey, err := hex.DecodeString(split[1])
	if err != nil {
		return []byte{}
	}
	return pubKey
}

func (t TransactionIn) GetAddress() wallet.Address {
	return wallet.GetAddress(t.GetPublicKey())
}

func (t *TransactionIn) GetPkHash() []byte {
	return t.GetAddress().GetScriptAddress()
}

func (t TransactionIn) GetAddressString() string {
	return t.GetAddress().GetEncoded()
}

func (t TransactionIn) Delete() error {
	if len(t.TxnOutHashString) != 0 {
		var txOut TransactionOut
		err := find(&txOut, TransactionOut{HashString: t.TxnOutHashString})
		if err != nil {
			if ! IsRecordNotFoundError(err) {
				return jerr.Get("error getting transaction out", err)
			}
		} else {
			txOut.TxnInHashString = ""
			err = txOut.Save()
			if err != nil {
				return jerr.Get("error saving transaction out", err)
			}
		}
	}
	result := remove(t)
	if result.Error != nil {
		return jerr.Get("error removing transaction input", result.Error)
	}
	return nil
}

func GetTransactionInputByHashString(hashString string) (*TransactionIn, error) {
	var txIn TransactionIn
	err := findPreloadColumns(transactionInColumns, &txIn, TransactionIn{
		HashString: hashString,
	})
	if err != nil {
		return nil, jerr.Get("error finding transaction in", err)
	}
	return &txIn, nil
}

func GetTransactionInput(txHash []byte, index uint32) (*TransactionIn, error) {
	var transactionIn TransactionIn
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	result := db.
	Where("`previous_out_point_index` = ?", index).
		Find(&transactionIn, TransactionIn{PreviousOutPointHash: txHash})
	if result.Error != nil {
		return nil, jerr.Get("error finding transaction input", result.Error)
	}
	return &transactionIn, nil
}
