package feed_event

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func AddPost(post *db.MemoPost) error {
	var feed = db.Feed{
		PkHash: post.PkHash,
		TxHash: post.TxHash,
	}
	if post.Block != nil {
		feed.BlockHeight = post.Block.Height
	}
	if post.ParentTxHash != nil {
		feed.EventType = db.FeedEventReply
	} else if post.Topic != "" {
		feed.EventType = db.FeedEventTopicPost
	} else if post.IsPoll {
		feed.EventType = db.FeedEventCreatePoll
	} else {
		feed.EventType = db.FeedEventPost
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed post", err)
	}
	return nil
}
