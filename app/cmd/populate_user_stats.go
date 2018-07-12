package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/spf13/cobra"
)

var populateUserStatsCmd = &cobra.Command{
	Use:  "populate-user-stats",
	RunE: func(c *cobra.Command, args []string) error {
		userStats, err := db.GetUserStats()
		if err != nil {
			jerr.Get("error getting user stats", err).Print()
			return nil
		}
		for _, userStat := range userStats {
			_, err := db.AddUpdateStat(userStat)
			if err != nil {
				jerr.Get("error updating user stat", err).Print()
				return nil
			}
		}
		fmt.Printf("Updated %d user stats\n", len(userStats))
		return nil
	},
}
