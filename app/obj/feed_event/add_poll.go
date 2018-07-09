package feed_event

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func AddPollVote(pollVote *db.MemoPollVote) error {
	var feed = db.FeedEvent{
		PkHash:    pollVote.PkHash,
		TxHash:    pollVote.TxHash,
		EventType: db.FeedEventPollVote,
	}
	if pollVote.Block != nil {
		feed.BlockHeight = pollVote.Block.Height
	}
	err := feed.Save()
	if err != nil {
		return jerr.Get("error saving feed poll vote", err)
	}
	return nil
}
