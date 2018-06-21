package feed_event

import (
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/util"
)

type Event struct {
	FeedEvent *db.FeedEvent
	Name      string
	Post      *profile.Post
	SetName   *db.MemoSetName
}

func (e *Event) TimeAgo() string {
	if e.FeedEvent.Block != nil && e.FeedEvent.Block.Timestamp.Before(e.FeedEvent.CreatedAt) {
		return util.GetTimeAgo(e.FeedEvent.Block.Timestamp)
	}
	return util.GetTimeAgo(e.FeedEvent.CreatedAt)
}

func (e *Event) IsLike() bool {
	return e.FeedEvent.EventType == db.FeedEventLike
}

func (e *Event) IsPost() bool {
	return e.FeedEvent.EventType == db.FeedEventPost
}

func (e *Event) IsReply() bool {
	return e.FeedEvent.EventType == db.FeedEventReply
}

func (e *Event) IsTopicPost() bool {
	return e.FeedEvent.EventType == db.FeedEventTopicPost
}

func (e *Event) IsCreatePoll() bool {
	return e.FeedEvent.EventType == db.FeedEventCreatePoll
}

func (e *Event) IsPollVote() bool {
	return e.FeedEvent.EventType == db.FeedEventPollVote
}

func (e *Event) IsSetName() bool {
	return e.FeedEvent.EventType == db.FeedEventSetName
}

func (e *Event) IsSetProfile() bool {
	return e.FeedEvent.EventType == db.FeedEventSetProfile
}

func (e *Event) IsSetProfilePic() bool {
	return e.FeedEvent.EventType == db.FeedEventSetProfilePic
}

func (e *Event) IsFollowUser() bool {
	return e.FeedEvent.EventType == db.FeedEventFollowUser
}

func (e *Event) IsFollowTopic() bool {
	return e.FeedEvent.EventType == db.FeedEventFollowTopic
}

func (e *Event) GetType() string {
	switch e.FeedEvent.EventType {
	case db.FeedEventPost:
		return "Post"
	case db.FeedEventReply:
		return "Reply"
	case db.FeedEventLike:
		return "Like"
	case db.FeedEventTopicPost:
		return "Topic Post"
	case db.FeedEventCreatePoll:
		return "Create Poll"
	case db.FeedEventPollVote:
		return "Poll Vote"
	case db.FeedEventSetName:
		return "Set Name"
	case db.FeedEventSetProfile:
		return "Set Profile"
	case db.FeedEventSetProfilePic:
		return "Set Profile Pic"
	case db.FeedEventFollowUser:
		return "Follow User"
	case db.FeedEventFollowTopic:
		return "Follow Topic"
	}
	return ""
}
