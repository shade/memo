package feed_event

import (
	"bytes"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
)

func GetEventsForUser(pkHash []byte, offset uint) ([]*Event, error) {
	var events []*Event
	feedEvents, err := db.GetRecentFeedForPkHash(pkHash, uint(offset))
	if err != nil {
		return nil, jerr.Get("error getting recent feed for pk hash", err)
	}
	var (
		postTxHashes          [][]byte
		likeTxHashes          [][]byte
		postVoteTxHashes      [][]byte
		setNameTxHashes       [][]byte
		profileSetTxHashes    [][]byte
		profileSetPicTxHashes [][]byte
		getPkHashNames        [][]byte
	)
	for _, feedEvent := range feedEvents {
		switch feedEvent.EventType {
		case db.FeedEventPost, db.FeedEventTopicPost, db.FeedEventReply, db.FeedEventCreatePoll:
			postTxHashes = append(postTxHashes, feedEvent.TxHash)
		case db.FeedEventLike:
			likeTxHashes = append(likeTxHashes, feedEvent.TxHash)
		case db.FeedEventPollVote:
			postVoteTxHashes = append(postVoteTxHashes, feedEvent.TxHash)
		case db.FeedEventSetName:
			setNameTxHashes = append(setNameTxHashes, feedEvent.TxHash)
		case db.FeedEventSetProfile:
			profileSetTxHashes = append(profileSetTxHashes, feedEvent.TxHash)
		case db.FeedEventSetProfilePic:
			profileSetPicTxHashes = append(profileSetPicTxHashes, feedEvent.TxHash)
		}
		getPkHashNames = append(getPkHashNames, feedEvent.PkHash)
	}

	feedEventMemoSetNames, err := db.GetNamesForPkHashes(getPkHashNames)
	if err != nil {
		return nil, jerr.Get("error getting feed event set names", err)
	}

	memoSetNames, err := db.GetSetNamesByTxHashes(setNameTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting set names by tx hashes", err)
	}

	memoLikes, err := db.GetMemoLikesByTxHashes(likeTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting memo likes by tx hashes", err)
	}
	for _, memoLike := range memoLikes {
		postTxHashes = append(postTxHashes, memoLike.LikeTxHash)
	}

	memoPollVotes, err := db.GetMemoPollVotesByTxHashes(postVoteTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting memo poll votes by tx hashes", err)
	}
	var optionTxHashes [][]byte
	for _, memoPollVote := range memoPollVotes {
		optionTxHashes = append(optionTxHashes, memoPollVote.OptionTxHash)
	}

	memoPollOptions, err := db.GetMemoPollOptionsByTxHashes(optionTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting memo poll options by tx hashes", err)
	}
	for _, memoPollOption := range memoPollOptions {
		postTxHashes = append(postTxHashes, memoPollOption.PollTxHash)
	}

	for i := 0; i < len(postTxHashes); i++ {
		for g := i + 1; g < len(postTxHashes); g++ {
			if bytes.Equal(postTxHashes[i], postTxHashes[g]) {
				postTxHashes = append(postTxHashes[:g], postTxHashes[g+1:]...)
				g--
			}
		}
	}

	posts, err := profile.GetPostsByTxHashes(postTxHashes, pkHash)
	if err != nil {
		return nil, jerr.Get("error getting memo posts by tx hashes", err)
	}
	for _, feedEvent := range feedEvents {
		var event = &Event{
			FeedEvent: feedEvent,
		}
		switch feedEvent.EventType {
		case db.FeedEventPost, db.FeedEventTopicPost, db.FeedEventReply, db.FeedEventCreatePoll:
			for _, post := range posts {
				if bytes.Equal(post.Memo.TxHash, feedEvent.TxHash) {
					event.Post = post
				}
			}
		case db.FeedEventLike:
			for _, memoLike := range memoLikes {
				if bytes.Equal(memoLike.TxHash, feedEvent.TxHash) {
					for _, post := range posts {
						if bytes.Equal(post.Memo.TxHash, memoLike.LikeTxHash) {
							event.Post = post
						}
					}
				}
			}
		case db.FeedEventPollVote:
			for _, memoPollVote := range memoPollVotes {
				if bytes.Equal(memoPollVote.TxHash, feedEvent.TxHash) {
					for _, memoPollOption := range memoPollOptions {
						if bytes.Equal(memoPollOption.TxHash, memoPollVote.OptionTxHash) {
							for _, post := range posts {
								if bytes.Equal(post.Memo.TxHash, memoPollOption.PollTxHash) {
									event.Post = post
								}
							}
						}
					}
				}
			}
		case db.FeedEventSetName:
			for _, memoSetName := range memoSetNames {
				if bytes.Equal(memoSetName.TxHash, feedEvent.TxHash) {
					event.SetName = memoSetName
				}
			}
		case db.FeedEventSetProfile:
		case db.FeedEventSetProfilePic:
		case db.FeedEventFollowUser:
		case db.FeedEventFollowTopic:
		default:
			jerr.Newf("unable to match feed event type: %#v", feedEvent).Print()
		}
		for _, feedEventMemoSetName := range feedEventMemoSetNames {
			if bytes.Equal(feedEventMemoSetName.PkHash, feedEvent.PkHash) {
				event.Name = feedEventMemoSetName.Name
			}
		}
		events = append(events, event)
	}
	return events, nil
}
