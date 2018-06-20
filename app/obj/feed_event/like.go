package feed_event

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func AddLike(like *db.MemoLike) error {
	var feed = db.Feed{
		PkHash:    like.PkHash,
		TxHash:    like.TxHash,
		EventType: db.FeedEventLike,
	}
	if like.Block != nil {
		feed.BlockHeight = like.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed like", err)
	}
	return nil
}
