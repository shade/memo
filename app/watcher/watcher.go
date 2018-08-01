package watcher

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"time"
)

type txSend struct {
	Hash string
	Type uint
}

func init() {
	var lastPing time.Time
	go func() {
		for {
			var needsPing bool
			if time.Since(lastPing) > 10 * time.Second {
				needsPing = true
				lastPing = time.Now()
			}
			var topicLastPostIds = make(map[string]uint)
			for i := 0; i < len(topicSockets); i++ {
				var topicSocket = topicSockets[i]
				_, ok := topicLastPostIds[topicSocket.Topic]
				if !ok {
					topicLastPostIds[topicSocket.Topic] = topicSocket.LastPostId
				}
				if topicSocket.LastPostId < topicLastPostIds[topicSocket.Topic] {
					topicLastPostIds[topicSocket.Topic] = topicSocket.LastPostId
				}
				if needsPing {
					err := topicSocket.Socket.Ping()
					if err != nil {
						go func(socket *TopicSocket, err error) {
							socket.Error <- nil
						}(topicSocket, err)
						topicSockets = append(topicSockets[:i], topicSockets[i+1:]...)
						i--
					}
				}
			}
			for topic, lastPostId := range topicLastPostIds {
				recentPosts, err := db.GetRecentPostsForTopic(topic, lastPostId)
				if err != nil && !db.IsRecordNotFoundError(err) {
					for i := 0; i < len(topicSockets); i++ {
						var topicSocket = topicSockets[i]
						if topicSocket.Topic == topic {
							go func(socket *TopicSocket, err error) {
								socket.Error <- jerr.Get("error getting recent post for topic", err)
							}(topicSocket, err)
							topicSockets = append(topicSockets[:i], topicSockets[i+1:]...)
							i--
						}
					}
				}
				if len(recentPosts) > 0 {
					for _, recentPost := range recentPosts {
						txHash := recentPost.GetTransactionHashString()
						for i := 0; i < len(topicSockets); i++ {
							var topicSocket = topicSockets[i]
							if topicSocket.Topic == topic && topicSocket.LastPostId < recentPost.Id {
								topicSocket.LastPostId = recentPost.Id
								err = topicSocket.Socket.WriteJSON(txSend{
									Hash: txHash,
									Type: 1,
								})
								if err != nil {
									go func(socket *TopicSocket, err error) {
										socket.Error <- nil
									}(topicSocket, err)
									topicSockets = append(topicSockets[:i], topicSockets[i+1:]...)
									i--
								}
							}
						}
					}
				}
			}
			var topicLastLikeIds = make(map[string]uint)
			for _, topicSocket := range topicSockets {
				_, ok := topicLastLikeIds[topicSocket.Topic]
				if !ok {
					topicLastLikeIds[topicSocket.Topic] = topicSocket.LastLikeId
				}
				if topicSocket.LastLikeId < topicLastLikeIds[topicSocket.Topic] {
					topicLastLikeIds[topicSocket.Topic] = topicSocket.LastLikeId
				}
			}
			for topic, lastLikeId := range topicLastLikeIds {
				recentLikes, err := db.GetRecentLikesForTopic(topic, lastLikeId)
				if err != nil && !db.IsRecordNotFoundError(err) {
					for i := 0; i < len(topicSockets); i++ {
						var topicSocket = topicSockets[i]
						if topicSocket.Topic == topic {
							go func(socket *TopicSocket, err error) {
								socket.Error <- jerr.Get("error getting recent like for topic", err)
							}(topicSocket, err)
							topicSockets = append(topicSockets[:i], topicSockets[i+1:]...)
							i--
						}
					}
				}
				if len(recentLikes) > 0 {
					for _, recentLike := range recentLikes {
						txHash := recentLike.GetLikeTransactionHashString()
						for i := 0; i < len(topicSockets); i++ {
							var topicSocket = topicSockets[i]
							if topicSocket.Topic == topic && topicSocket.LastLikeId < recentLike.Id {
								topicSocket.LastLikeId = recentLike.Id
								err = topicSocket.Socket.WriteJSON(txSend{
									Hash: txHash,
									Type: 2,
								})
								if err != nil {
									go func(socket *TopicSocket, err error) {
										socket.Error <- jerr.Get("error writing to socket", err)
									}(topicSocket, err)
									topicSockets = append(topicSockets[:i], topicSockets[i+1:]...)
									i--
								}
							}
						}
					}
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
}
