package obj

import "time"

type UserStat struct {
	PkHash    []byte
	NumPosts  int
	FirstPost time.Time
	LastPost  time.Time
}
