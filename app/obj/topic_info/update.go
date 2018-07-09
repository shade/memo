package topic_info

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func Update(topicNames ...string) error {
	topics, err := db.GetTopicInfoFromPosts(topicNames...)
	if err != nil {
		return jerr.Get("error getting all unique topics", err)
	}
	for _, topic := range topics {
		err = db.AddUpdateTopicInfo(topic.Name, topic.CountPosts, topic.CountFollows, topic.RecentTime)
		if err != nil {
			return jerr.Get("error updating topic info", err)
		}
	}
	return nil
}
