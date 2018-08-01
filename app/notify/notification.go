package notify

import (
	"bytes"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/util"
	"time"
)

type NotificationType string

const (
	TypeLike   = "like"
	TypeFollow = "follow"
	TypeReply  = "reply"
)

type Notification struct {
	DbId           uint
	DbNotification *db.Notification
	Type           NotificationType
	Name           string
	PkHash         []byte
	AddressString  string
	PostHashString string
	Message        string
	Time           time.Time
	ProfilePic     *db.MemoSetPic
	// Reply
	ParentMessage    string
	ParentHashString string
	// Like
	TipAmount int64
}

func (n Notification) IsLike() bool {
	return n.Type == TypeLike
}

func (n Notification) IsReply() bool {
	return n.Type == TypeReply
}

func (n Notification) IsNewFollower() bool {
	return n.Type == TypeFollow
}

func (n Notification) GetTimeAgo() string {
	return util.GetTimeAgo(n.Time)
}

func (n Notification) GetId() uint {
	if n.DbNotification == nil {
		return 0
	}
	return n.DbNotification.Id
}

func AttachProfilePicsToNotifications(notifications []*Notification) error {
	var picPkHashes [][]byte
	for _, notification := range notifications {
		for _, namePkHash := range picPkHashes {
			if bytes.Equal(namePkHash, notification.PkHash) {
				continue
			}
		}
		picPkHashes = append(picPkHashes, notification.PkHash)
	}
	setPics, err := db.GetPicsForPkHashes(picPkHashes)
	if err != nil {
		return jerr.Get("error getting profile pics for pk hashes", err)
	}
	for _, setPic := range setPics {
		for _, notification := range notifications {
			if bytes.Equal(notification.PkHash, setPic.PkHash) {
				notification.ProfilePic = setPic
			}
		}
	}
	return nil
}

func AttachNamesToNotifications(notifications []*Notification) error {
	var namePkHashes [][]byte
	for _, notification := range notifications {
		for _, namePkHash := range namePkHashes {
			if bytes.Equal(namePkHash, notification.PkHash) {
				continue
			}
		}
		namePkHashes = append(namePkHashes, notification.PkHash)
	}
	setNames, err := db.GetNamesForPkHashes(namePkHashes)
	if err != nil {
		return jerr.Get("error getting set names for pk hashes", err)
	}
	for _, notification := range notifications {
		for _, setName := range setNames {
			if bytes.Equal(notification.PkHash, setName.PkHash) {
				notification.Name = setName.Name
			}
		}
		if notification.Name == "" {
			notification.Name = notification.AddressString[:16]
		}
	}
	return nil
}
