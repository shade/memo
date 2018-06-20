package feed_event

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func AddSetName(setName *db.MemoSetName) error {
	var feed = db.Feed{
		PkHash:    setName.PkHash,
		TxHash:    setName.TxHash,
		EventType: db.FeedEventSetName,
	}
	if setName.Block != nil {
		feed.BlockHeight = setName.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed set name", err)
	}
	return nil
}

func AddSetProfile(setProfile *db.MemoSetProfile) error {
	var feed = db.Feed{
		PkHash:    setProfile.PkHash,
		TxHash:    setProfile.TxHash,
		EventType: db.FeedEventSetProfile,
	}
	if setProfile.Block != nil {
		feed.BlockHeight = setProfile.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed set profile", err)
	}
	return nil
}

func AddSetProfilePic(setName *db.MemoSetPic) error {
	var feed = db.Feed{
		PkHash:    setName.PkHash,
		TxHash:    setName.TxHash,
		EventType: db.FeedEventSetProfilePic,
	}
	if setName.Block != nil {
		feed.BlockHeight = setName.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed set profile pic", err)
	}
	return nil
}
