package view

import (
	"github.com/memocash/memo/app/util"
	"net/url"
	"time"
)

type Topic struct {
	Name         string
	RecentTime   time.Time
	CountPosts   int
	CountFollows int
	UnreadPosts  bool
}

func (t Topic) GetUrlEncoded() string {
	return url.QueryEscape(t.Name)
}

func (t Topic) GetTimeAgo() string {
	return util.GetTimeAgo(t.RecentTime)
}

type TopicOrderType int

const (
	TopicOrderTypeRecent    TopicOrderType = iota
	TopicOrderTypeFollowers
	TopicOrderTypePosts
)
