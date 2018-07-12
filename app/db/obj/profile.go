package obj

import "time"

type Profile struct {
	PkHash       []byte
	NumPosts     int
	FirstPost    time.Time
	LastPost     time.Time
	NumFollowers int
}
