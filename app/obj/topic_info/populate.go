package topic_info

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func PopulateAll() error {
	fmt.Println("Querying topic info from posts...")
	topics, err := db.GetTopicInfoFromPosts()
	if err != nil {
		return jerr.Get("error getting all unique topics", err)
	}
	fmt.Printf("Updating %d topics\n", len(topics))
	for _, topic := range topics {
		err = db.AddUpdateTopicInfo(topic.Name, topic.CountPosts, topic.CountFollows, topic.RecentTime)
		if err != nil {
			return jerr.Get("error updating topic info", err)
		}
	}
	fmt.Println("All done.")
	return nil
}
