package cmd

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/obj/topic_info"
	"github.com/spf13/cobra"
)

var populateTopicInfoCmd = &cobra.Command{
	Use:  "populate-topic-info",
	RunE: func(c *cobra.Command, args []string) error {
		err := topic_info.PopulateAll()
		if err != nil {
			jerr.Get("error populating all topics", err).Print()
		}
		return nil
	},
}
