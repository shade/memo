package feed_event

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func AddFollow(follow *db.MemoFollow) error {
	var feed = db.Feed{
		PkHash:    follow.PkHash,
		TxHash:    follow.TxHash,
		EventType: db.FeedEventFollowUser,
	}
	if follow.Block != nil {
		feed.BlockHeight = follow.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed follow", err)
	}
	return nil
}
