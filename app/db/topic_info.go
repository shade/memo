package db

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db/view"
	"time"
)

type TopicInfo struct {
	Id            uint `gorm:"primary_key"`
	TopicName     string
	PostCount     int
	FollowerCount int
	RecentPost    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func AddUpdateTopicInfo(topicName string, numPosts int, followers int, recentPost time.Time) error {
	var topicInfo = TopicInfo{
		TopicName: topicName,
	}
	err := find(&topicInfo, topicInfo)
	if err != nil && ! IsRecordNotFoundError(err) {
		return jerr.Get("error getting topic info", err)
	}
	topicInfo.PostCount = numPosts
	topicInfo.FollowerCount = followers
	topicInfo.RecentPost = recentPost
	result := save(&topicInfo)
	if result.Error != nil {
		return jerr.Get("error saving topic info", result.Error)
	}
	return nil
}

func GetTopicInfo(offset uint, searchString string, pkHash []byte, orderType view.TopicOrderType) ([]*view.Topic, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Limit(25).
		Offset(offset)
	if searchString != "" {
		query = query.Where("topic_name LIKE ?", fmt.Sprintf("%%%s%%", searchString))
	}
	if len(pkHash) > 0 {
		joinQuery := "JOIN (" +
			"SELECT MAX(id) AS id " +
			"FROM memo_topic_follows " +
			"WHERE pk_hash = ? " +
			"GROUP BY topic" +
			") sq ON (sq.id = memo_topic_follows.id)"
		query = query.
			Joins("JOIN memo_topic_follows ON (topic_infos.topic_name = memo_topic_follows.topic)").
			Joins(joinQuery, pkHash).
			Where("memo_topic_follows.unfollow = 0")
	}
	switch orderType {
	case view.TopicOrderTypeFollowers:
		query = query.Order("follower_count DESC")
	case view.TopicOrderTypePosts:
		query = query.Order("post_count DESC")
	}
	query = query.Order("recent_post DESC")
	var topicInfos []*TopicInfo
	result := query.Find(&topicInfos)
	if result.Error != nil {
		return nil, jerr.Get("error getting topic infos", result.Error)
	}
	var topics []*view.Topic
	for _, topicInfo := range topicInfos {
		topics = append(topics, &view.Topic{
			Name:         topicInfo.TopicName,
			RecentTime:   topicInfo.RecentPost,
			CountPosts:   topicInfo.PostCount,
			CountFollows: topicInfo.FollowerCount,
		})
	}
	return topics, nil
}
