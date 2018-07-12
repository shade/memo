package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/obj/user_stats"
	"github.com/spf13/cobra"
)

var populateUserStatsCmd = &cobra.Command{
	Use:  "populate-user-stats",
	RunE: func(c *cobra.Command, args []string) error {
		err := user_stats.Populate()
		if err != nil {
			jerr.Get("error populating user stats", err).Print()
			return nil
		}
		fmt.Println("All done.")
		return nil
	},
}
