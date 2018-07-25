package feed_event

import (
	"bytes"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/rep"
	"github.com/memocash/memo/app/profile"
)

func GetEventsForUser(userId uint, pkHash []byte, offset uint) ([]*Event, error) {
	feedEvents, err := db.GetRecentFeedForPkHash(pkHash, uint(offset))
	if err != nil {
		return nil, jerr.Get("error getting recent feed for pk hash", err)
	}
	events, err := getEvents(feedEvents, userId, pkHash)
	if err != nil {
		return nil, jerr.Get("error getting events from feed events", err)
	}
	return events, nil
}

func GetUserEvents(userId uint, userPkHash []byte, pkHash []byte, offset uint, eventTypes []db.FeedEventType) ([]*Event, error) {
	feedEvents, err := db.GetRecentFeedUserEvents(pkHash, uint(offset), eventTypes)
	if err != nil {
		return nil, jerr.Get("error getting recent user feed for pk hash", err)
	}
	events, err := getEvents(feedEvents, userId, userPkHash)
	if err != nil {
		return nil, jerr.Get("error getting events from feed events", err)
	}
	return events, nil
}

func GetAllEvents(userId uint, pkHash []byte, offset uint) ([]*Event, error) {
	feedEvents, err := db.GetRecentFeedEvents(uint(offset))
	if err != nil {
		return nil, jerr.Get("error getting recent feed events", err)
	}
	events, err := getEvents(feedEvents, userId, pkHash)
	if err != nil {
		return nil, jerr.Get("error getting events from feed events", err)
	}
	return events, nil
}

func getEvents(feedEvents []*db.FeedEvent, userId uint, pkHash []byte) ([]*Event, error) {
	var (
		postTxHashes          [][]byte
		likeTxHashes          [][]byte
		postVoteTxHashes      [][]byte
		setNameTxHashes       [][]byte
		profileSetTxHashes    [][]byte
		profileSetPicTxHashes [][]byte
		userFollowTxHashes    [][]byte
		topicFollowTxHashes   [][]byte
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
		case db.FeedEventFollowUser:
			userFollowTxHashes = append(userFollowTxHashes, feedEvent.TxHash)
		case db.FeedEventFollowTopic:
			topicFollowTxHashes = append(topicFollowTxHashes, feedEvent.TxHash)
		}
		getPkHashNames = append(getPkHashNames, feedEvent.PkHash)
	}

	memoFollows, err := db.GetMemoFollowsByTxHashes(userFollowTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting follows by tx hashes", err)
	}

	for _, memoFollow := range memoFollows {
		getPkHashNames = append(getPkHashNames, memoFollow.FollowPkHash)
	}

	feedEventMemoSetNames, err := db.GetNamesForPkHashes(getPkHashNames)
	if err != nil {
		return nil, jerr.Get("error getting feed event set names", err)
	}

	feedEventMemoSetProfilePics, err := db.GetPicsForPkHashes(getPkHashNames)
	if err != nil {
		return nil, jerr.Get("error getting feed event set profile pics", err)
	}

	memoSetNames, err := db.GetSetNamesByTxHashes(setNameTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting set names by tx hashes", err)
	}

	memoSetProfiles, err := db.GetSetProfilesByTxHashes(profileSetTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting set profiles by tx hashes", err)
	}

	memoSetProfilePics, err := db.GetSetProfilePicsByTxHashes(profileSetPicTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting set profile pics by tx hashes", err)
	}

	memoTopicFollows, err := db.GetMemoTopicFollowsByTxHashes(topicFollowTxHashes)
	if err != nil {
		return nil, jerr.Get("error getting topic follows by tx hashes", err)
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
	err = profile.SetShowMediaForPosts(posts, userId)
	if err != nil {
		return nil, jerr.Get("error setting show media for posts", err)
	}
	err = profile.AttachPollsToPosts(posts)
	if err != nil {
		return nil, jerr.Get("error attaching polls to posts", err)
	}
	err = profile.AttachParentToPosts(posts)
	if err != nil {
		return nil, jerr.Get("error attaching parents to posts", err)
	}
	err = profile.AttachLikesToPosts(posts)
	if err != nil {
		return nil, jerr.Get("error attaching parents to posts", err)
	}
	err = profile.AttachProfilePicsToPosts(posts)
	if err != nil {
		return nil, jerr.Get("error attaching profile pics to posts", err)
	}
	err = profile.AttachReputationToPosts(posts)
	if err != nil {
		return nil, jerr.Get("error attaching profile pics to posts", err)
	}
	settings, err := cache.GetUserSettings(userId)
	if err != nil {
		return nil, jerr.Get("error getting user settings", err)
	}
	var showMedia bool
	if settings.Integrations == db.SettingIntegrationsAll {
		showMedia = true
		for _, post := range posts {
			post.ShowMedia = true
			if post.Parent != nil {
				post.Parent.ShowMedia = true
			}
		}
	}
	var events []*Event
	for _, feedEvent := range feedEvents {
		var event = &Event{
			FeedEvent:  feedEvent,
			SelfPkHash: pkHash,
			ShowMedia:  showMedia,
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
					event.MemoLike = memoLike
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
					event.PollVote = memoPollVote
					for _, memoPollOption := range memoPollOptions {
						if bytes.Equal(memoPollOption.TxHash, memoPollVote.OptionTxHash) {
							event.PollOption = memoPollOption
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
			for _, memoSetProfile := range memoSetProfiles {
				if bytes.Equal(memoSetProfile.TxHash, feedEvent.TxHash) {
					event.SetProfile = memoSetProfile
				}
			}
		case db.FeedEventSetProfilePic:
			for _, memoSetPic := range memoSetProfilePics {
				if bytes.Equal(memoSetPic.TxHash, feedEvent.TxHash) {
					event.SetProfilePic = memoSetPic
				}
			}
		case db.FeedEventFollowUser:
			for _, memoFollow := range memoFollows {
				if bytes.Equal(memoFollow.TxHash, feedEvent.TxHash) {
					event.UserFollow = memoFollow
					for _, memoSetName := range feedEventMemoSetNames {
						if bytes.Equal(memoSetName.PkHash, memoFollow.FollowPkHash) {
							event.FollowName = memoSetName.Name
						}
					}
					for _, memoSetPic := range feedEventMemoSetProfilePics {
						if bytes.Equal(memoSetPic.PkHash, memoFollow.FollowPkHash) {
							event.FollowProfilePic = memoSetPic
						}
					}
				}
			}
		case db.FeedEventFollowTopic:
			for _, memoTopicFollow := range memoTopicFollows {
				if bytes.Equal(memoTopicFollow.TxHash, feedEvent.TxHash) {
					event.TopicFollow = memoTopicFollow
				}
			}
		default:
			jerr.Newf("unable to match feed event type: %#v", feedEvent).Print()
		}
		for _, feedEventMemoSetName := range feedEventMemoSetNames {
			if bytes.Equal(feedEventMemoSetName.PkHash, feedEvent.PkHash) {
				event.Name = feedEventMemoSetName.Name
			}
		}
		for _, feedEventMemoSetProfilePic := range feedEventMemoSetProfilePics {
			if bytes.Equal(feedEventMemoSetProfilePic.PkHash, feedEvent.PkHash) {
				event.ProfilePic = feedEventMemoSetProfilePic
			}
		}
		events = append(events, event)
	}
	err = AttachReputationToEvents(events)
	if err != nil {
		jerr.Get("error attaching reputation to events", err).Print()
	}
	return events, nil
}

func AttachReputationToEvents(events []*Event) error {
	for _, event := range events {
		reputation, err := rep.GetReputation(event.SelfPkHash, event.FeedEvent.PkHash)
		if err != nil {
			return jerr.Get("error getting reputation", err)
		}
		event.Reputation = reputation
		if event.FeedEvent.EventType == db.FeedEventFollowUser {
			followReputation, err := rep.GetReputation(event.SelfPkHash, event.UserFollow.FollowPkHash)
			if err != nil {
				return jerr.Get("error getting follow reputation", err)
			}
			event.FollowReputation = followReputation
		}
	}
	return nil
}
