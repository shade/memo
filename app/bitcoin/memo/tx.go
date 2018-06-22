package memo

import (
	"github.com/jchavannes/btcd/wire"
	"github.com/memocash/memo/app/bitcoin/wallet"
)

type SpendOutputType uint

const (
	SpendOutputTypeP2PK                   SpendOutputType = iota
	SpendOutputTypeReturn
	SpendOutputTypeMemoMessage
	SpendOutputTypeMemoSetName
	SpendOutputTypeMemoFollow
	SpendOutputTypeMemoUnfollow
	SpendOutputTypeMemoLike
	SpendOutputTypeMemoReply
	SpendOutputTypeMemoSetProfile
	SpendOutputTypeMemoTopicMessage
	SpendOutputTypeMemoTopicFollow
	SpendOutputTypeMemoTopicUnfollow
	SpendOutputTypeMemoPollQuestionSingle
	SpendOutputTypeMemoPollQuestionMulti
	SpendOutputTypeMemoPollOption
	SpendOutputTypeMemoPollVote
	SpendOutputTypeMemoSetProfilePic
)

type SpendOutput struct {
	Address wallet.Address
	Amount  int64
	Type    SpendOutputType
	RefData []byte
	Data    []byte
}

type TxInput struct {
	PkHash      []byte
	Value       int64
	PrevOutHash string
}

type Tx struct {
	SelfPkHash []byte
	MsgTx      *wire.MsgTx
	Inputs     []*TxInput
}
