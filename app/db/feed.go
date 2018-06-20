package db

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type FeedEvent int

const (
	FeedEventPost          FeedEvent = 1
	FeedEventReply         FeedEvent = 2
	FeedEventLike          FeedEvent = 3
	FeedEventTopicPost     FeedEvent = 4
	FeedEventCreatePoll    FeedEvent = 5
	FeedEventPollVote      FeedEvent = 6
	FeedEventSetName       FeedEvent = 7
	FeedEventSetProfile    FeedEvent = 8
	FeedEventSetProfilePic FeedEvent = 9
	FeedEventFollowUser    FeedEvent = 10
	FeedEventFollowTopic   FeedEvent = 11
)

type Feed struct {
	Id          uint   `gorm:"primary_key"`
	BlockHeight uint   `gorm:"index:block_height"`
	Block       *Block
	PkHash      []byte `gorm:"index:pk_hash"`
	TxHash      []byte `gorm:"unique"`
	EventType   FeedEvent
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (f *Feed) Save() error {
	result := save(&f)
	if result.Error == nil {
		return nil
	} else if !IsAlreadyExistsError(result.Error) {
		return jerr.Get("error saving feed event", result.Error)
	}
	if f.BlockHeight == 0 {
		return nil
	}
	feed, err := GetFeedByTxHash(f.TxHash)
	if err != nil {
		return jerr.Get("error getting feed event", result.Error)
	}
	if feed.BlockHeight != 0 {
		return nil
	}
	feed.BlockHeight = f.BlockHeight
	result = save(&feed)
	if result.Error != nil {
		return jerr.Get("error updating feed event", result.Error)
	}
	return nil
}

func GetFeedByTxHash(txHash []byte) (*Feed, error) {
	var feed Feed
	err := find(feed, Feed{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting feed item", err)
	}
	return &feed, nil
}
