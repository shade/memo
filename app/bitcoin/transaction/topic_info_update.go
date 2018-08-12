package transaction

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/obj/topic_info"
)

func updateTopicInfo(topicName string) {
	go func() {
		err := topic_info.Update(topicName)
		if err != nil {
			jerr.Getf(err, "error updating topic info: %s", topicName).Print()
		}
	}()
}
