package view

import (
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/util"
	"net/url"
	"time"
)

type Thread struct {
	Topic       string
	Message     string
	RootTxHash  []byte
	NumReplies  int
	RecentReply time.Time
}

func (t Thread) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(t.RootTxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (t Thread) GetRecentTimeAgo() string {
	return util.GetTimeAgo(t.RecentReply)
}

func (t Thread) GetUrlEncodedTopic() string {
	return url.QueryEscape(t.Topic)
}
