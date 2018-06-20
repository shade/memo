package feed_event

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func AddTopicFollow(topicFollow *db.MemoTopicFollow) error {
	var feed = db.Feed{
		PkHash:    topicFollow.PkHash,
		TxHash:    topicFollow.TxHash,
		EventType: db.FeedEventFollowTopic,
	}
	if topicFollow.Block != nil {
		feed.BlockHeight = topicFollow.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed topic follow", err)
	}
	return nil
}
