package util

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func GetTimeAgo(ts time.Time) string {
	delta := time.Now().Sub(ts)
	hours := int(delta.Hours())
	if hours > 0 {
		if hours >= 24 {
			if hours < 48 {
				return "1 day ago"
			}
			return fmt.Sprintf("%d days ago", hours/24)
		}
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	minutes := int(delta.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	return fmt.Sprintf("%d seconds ago", int(delta.Seconds()))
}

func GetTimezoneTime(ts time.Time, timezone string) string {
	timeLayout := "2006-01-02 15:04:05"
	if len(timezone) > 0 {
		displayLocation, err := time.LoadLocation(timezone)
		if err != nil {
			jerr.Get("error finding location", err).Print()
			return ts.Format(timeLayout)
		}
		return ts.In(displayLocation).Format(timeLayout)
	} else {
		return ts.Format(timeLayout)
	}
}
