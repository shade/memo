package user_stats

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func Populate() error {
	userStats, err := db.GetUserStats()
	if err != nil {
		return jerr.Get("error getting user stats", err)
	}
	for _, userStat := range userStats {
		_, err := db.AddUpdateStat(userStat)
		if err != nil {
			return jerr.Get("error updating user stat", err)
		}
	}
	return nil
}
