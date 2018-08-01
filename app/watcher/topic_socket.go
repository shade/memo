package watcher

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/html-parser"
)

type TopicSocket struct {
	Socket     *web.Socket
	Topic      string
	LastPostId uint
	LastLikeId uint
	Error      chan error
}

var topicSockets []*TopicSocket

func RegisterSocket(socket *web.Socket, topic string, lastPostId uint, lastLikeId uint) error {
	topic = html_parser.EscapeWithEmojis(topic)
	var topicSocket = &TopicSocket{
		Socket:     socket,
		Topic:      topic,
		LastPostId: lastPostId,
		LastLikeId: lastLikeId,
		Error:      make(chan error),
	}
	topicSockets = append(topicSockets, topicSocket)
	return <-topicSocket.Error
}
