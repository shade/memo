package db

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db/obj"
	"time"
)

type UserStat struct {
	Id           uint   `gorm:"primary_key"`
	PkHash       []byte `gorm:"not null;unique"`
	NumPosts     int
	NumFollowers int
	FirstPost    time.Time
	LastPost     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *UserStat) Save() error {
	result := save(u)
	if result.Error != nil {
		return jerr.Get("error saving user stat", result.Error)
	}
	return nil
}

func GetUserStat(pkHash []byte) (*UserStat, error) {
	var userStat UserStat
	err := find(&userStat, UserStat{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting user stat", err)
	}
	return &userStat, nil
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
	userStat.NumFollowers = userStatObj.NumFollowers
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

type UserStatOrderType int

const (
	UserStatOrderCreated UserStatOrderType = iota
	UserStatOrderNewest
	UserStatOrderPosts
	UserStatOrderFollowers
)

func GetUniqueMemoAPkHashes(offset int, searchString string, orderType UserStatOrderType) ([]*obj.Profile, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var profiles []*obj.Profile
	query := db.
		Table("user_stats").
		Select("user_stats.pk_hash AS pk_hash," +
		"user_stats.num_followers AS num_followers," +
		"user_stats.num_posts AS num_posts," +
		"user_stats.first_post AS first_post," +
		"user_stats.last_post AS last_post").
		Limit(25).
		Offset(offset)
	if searchString != "" {
		joinSelect := "JOIN (" +
			"	SELECT MAX(id) AS id" +
			"	FROM memo_set_names" +
			"	GROUP BY pk_hash" +
			") sq ON (sq.id = memo_set_names.id)"
		query = query.
			Joins("JOIN memo_set_names ON (memo_set_names.pk_hash = user_stats.pk_hash)").
			Joins(joinSelect).
			Where("memo_set_names.name LIKE ?", fmt.Sprintf("%%%s%%", searchString))
	}
	if orderType == UserStatOrderPosts {
		query = query.Order("user_stats.num_posts DESC")
	} else if orderType == UserStatOrderFollowers {
		query = query.Order("user_stats.num_followers DESC")
	} else if orderType == UserStatOrderNewest {
		query = query.Order("user_stats.first_post DESC")
	} else {
		query = query.Order("user_stats.first_post ASC")
	}
	result := query.Find(&profiles)
	if result.Error != nil {
		return nil, jerr.Get("error getting profiles", result.Error)
	}
	return profiles, nil
}

func GetUniqueUserCount() (int, error) {
	cnt, err := count(UserStat{})
	if err != nil {
		return 0, jerr.Get("error getting user stat count", err)
	}
	return int(cnt), nil
}
