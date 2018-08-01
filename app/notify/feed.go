package notify

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func GetNotificationsFeed(pkHash []byte, offset uint) ([]*Notification, error) {
	dbNotifications, err := db.GetRecentNotificationsForUser(pkHash, offset)
	if err != nil {
		return nil, jerr.Get("error getting notifications from db", err)
	}
	var notifications []*Notification
	for _, dbNotification := range dbNotifications {
		switch dbNotification.Type {
		case db.NotificationTypeLike:
			like, err := db.GetMemoLike(dbNotification.TxHash)
			if err != nil {
				jerr.Get("error getting notification like", err).Print()
				continue
			}
			post, err := db.GetMemoPost(like.LikeTxHash)
			if err != nil {
				jerr.Get("error getting like post for notification", err).Print()
				continue
			}
			notifications = append(notifications, LikeNotification{
				Notification: dbNotification,
				Like:         like,
				Post:         post,
			}.GetNotification())
		case db.NotificationTypeReply:
			post, err := db.GetMemoPost(dbNotification.TxHash)
			if err != nil {
				jerr.Get("error getting notification post", err).Print()
				continue
			}
			parent, err := db.GetMemoPost(post.ParentTxHash)
			if err != nil {
				jerr.Get("error getting notification post parent", err).Print()
				continue
			}
			notifications = append(notifications, ReplyNotification{
				Notification: dbNotification,
				Post:         post,
				Parent:       parent,
			}.GetNotification())
		case db.NotificationTypeNewFollower:
			follow, err := db.GetMemoFollow(dbNotification.TxHash)
			if err != nil {
				jerr.Get("error getting notification new follower", err).Print()
				continue
			}
			notifications = append(notifications, NewFollowerNotification{
				Notification: dbNotification,
				Follow:       follow,
			}.GetNotification())
		}
	}
	err = AttachNamesToNotifications(notifications)
	if err != nil {
		return nil, jerr.Get("error attaching names to notifications", err)
	}
	err = AttachProfilePicsToNotifications(notifications)
	if err != nil {
		return nil, jerr.Get("error attaching profile pics to notifications", err)
	}
	return notifications, nil
}
