package db

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db/obj"
	"time"
)

type UserStat struct {
	Id        uint   `gorm:"primary_key"`
	PkHash    []byte `gorm:"not null;unique"`
	NumPosts  int
	FirstPost time.Time
	LastPost  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserStat) Save() error {
	result := save(u)
	if result.Error != nil {
		return jerr.Get("error saving user stat", result.Error)
	}
	return nil
}

func AddUpdateStat(userStatObj obj.UserStat) (*UserStat, error) {
	var userStat = UserStat{
		PkHash: userStatObj.PkHash,
	}
	err := find(&userStat, userStat)
	if err != nil && ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting existing user stat from db", err)
	}
	userStat.NumPosts = userStatObj.NumPosts
	userStat.FirstPost = userStatObj.FirstPost
	userStat.LastPost = userStatObj.LastPost
	err = userStat.Save()
	if err != nil {
		return nil, jerr.Get("error saving user stat", err)
	}
	return &userStat, nil
}

func GetUserFirstPostStats() ([]obj.UserDateStat, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("user_stats").
		Select("DATE(first_post) AS date, " +
		"COUNT(*) AS num_users").
		Where("num_posts > 1").
		Group("date").
		Order("date ASC")
	var userDataStats []obj.UserDateStat
	result := query.Find(&userDataStats)
	if result.Error != nil {
		return nil, jerr.Get("error getting user first post stats", result.Error)
	}
	return userDataStats, nil
}

func GetUserLastPostStats() ([]obj.UserDateStat, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("user_stats").
		Select("DATE(last_post) AS date, " +
		"COUNT(*) AS num_users").
		Where("num_posts > 1").
		Group("date").
		Order("date ASC")
	var userDataStats []obj.UserDateStat
	result := query.Find(&userDataStats)
	if result.Error != nil {
		return nil, jerr.Get("error getting user last post stats", result.Error)
	}
	return userDataStats, nil
}
