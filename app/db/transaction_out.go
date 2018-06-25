package db

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"html"
	"strings"
	"time"
)

type TransactionOut struct {
	Id              uint           `gorm:"primary_key"`
	Index           uint32         `gorm:"unique_index:transaction_out_index;"`
	HashString      string
	TransactionHash []byte         `gorm:"unique_index:transaction_out_index;"`
	Transaction     *Transaction   `gorm:"foreignkey:TransactionHash"`
	KeyPkHash       []byte         `gorm:"index:pk_hash"`
	Key             *Key           `gorm:"foreignkey:KeyPkHash"`
	Value           int64
	PkScript        []byte
	LockString      string
	RequiredSigs    uint
	ScriptClass     uint
	TxnInHashString string         `gorm:"index:txn_in_hash_string"`
	TxnIn           *TransactionIn `gorm:"foreignkey:TxnInHashString"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func getHashString(txHash []byte, index uint32) string {
	hash, err := chainhash.NewHash(txHash)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("out:%s:%d", hash.String(), index)
}

func (t TransactionOut) GetHashString() string {
	return getHashString(t.TransactionHash, t.Index)
}

func (t TransactionOut) IsMemo() bool {
	return strings.HasPrefix(t.LockString, "OP_RETURN 6d")
}

func (t TransactionOut) GetPkHash() []byte {
	split := strings.Split(t.LockString, " ")
	if len(split) != 5 {
		return []byte{}
	}
	pubKey, err := hex.DecodeString(split[2])
	if err != nil {
		return []byte{}
	}
	return pubKey
}

func (t TransactionOut) GetAddressString() string {
	addressPkHash, err := btcutil.NewAddressPubKeyHash(t.KeyPkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error parsing address", err).Print()
		return ""
	}
	return addressPkHash.String()
}

func (t TransactionOut) Save() error {
	result := save(&t)
	if result.Error != nil {
		return jerr.Get("error saving transaction output", result.Error)
	}
	return nil
}

func (t TransactionOut) ValueInBCH() float64 {
	return float64(t.Value) * 1.e-8
}

func (t TransactionOut) HasIn() bool {
	return len(t.TxnInHashString) > 0
}

func (t TransactionOut) IsSpendable() bool {
	if len(t.TxnInHashString) > 0 {
		txIn, _ := GetTransactionInputByHashString(t.TxnInHashString)
		if txIn.Transaction.BlockId > 0 {
			return false
		}
	}
	return true
}

func (t TransactionOut) GetScriptClass() string {
	return txscript.ScriptClass(t.ScriptClass).String()
}

func (t TransactionOut) GetMessage() string {
	if txscript.ScriptClass(t.ScriptClass) == txscript.NullDataTy {
		data, err := txscript.PushedData(t.PkScript)
		if err != nil || len(data) == 0 {
			return ""
		}
		return string(data[0])
	}
	return html.EscapeString(script.GetScriptString(t.PkScript))
}

func GetTransactionOutput(txHash []byte, index uint32) (*TransactionOut, error) {
	var transactionOut TransactionOut
	err := find(&transactionOut, TransactionOut{
		TransactionHash: txHash,
		Index:           index,
	})
	if err != nil {
		return nil, jerr.Get("error finding transaction output", err)
	}
	return &transactionOut, nil
}

func GetSpendableTransactionOutputsForPkHash(pkHash []byte) ([]*TransactionOut, error) {
	var transactionOuts []*TransactionOut
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	result := db.
		Where("txn_in_hash_string = ''").
		Where("value > 0").
		Find(&transactionOuts, TransactionOut{KeyPkHash: pkHash})
	if result.Error != nil {
		return nil, jerr.Get("error getting transaction outputs", result.Error)
	}
	return transactionOuts, nil
}

type TxOutSortByValue []*TransactionOut

func (txOuts TxOutSortByValue) Len() int      { return len(txOuts) }
func (txOuts TxOutSortByValue) Swap(i, j int) { txOuts[i], txOuts[j] = txOuts[j], txOuts[i] }
func (txOuts TxOutSortByValue) Less(i, j int) bool {
	return txOuts[i].Value > txOuts[j].Value
}

func HasSpendable(pkHash []byte) (bool, error) {
	transactionOutputs, err := GetSpendableTransactionOutputsForPkHash(pkHash)
	if err != nil {
		return false, jerr.Get("error getting transactions", err)
	}
	var totalValue int64
	for _, transactionOutput := range transactionOutputs {
		totalValue += transactionOutput.Value
	}
	return totalValue > 1000, nil
}
