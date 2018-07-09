package memo

import (
	"github.com/jchavannes/btcd/wire"
)

type TxInput struct {
	PkHash      []byte
	Value       int64
	PrevOutHash string
}

type Tx struct {
	Type       OutputType
	SelfPkHash []byte
	MsgTx      *wire.MsgTx
	Inputs     []*TxInput
}
