package view

import "time"

type UserStat struct {
	PkHash       []byte
	NumPosts     int
	NumFollowers int
	FirstPost    time.Time
	LastPost     time.Time
}
