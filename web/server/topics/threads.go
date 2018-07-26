package topics

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/res"
	"net/http"
	"net/url"
	"strings"
)

var threadsRoute = web.Route{
	Pattern: res.UrlTopicThreads + "/" + urlTopicName.UrlPart(),
	Handler: func(r *web.Response) {
		preHandler(r)
		topicRaw := r.Request.GetUrlNamedQueryVariable(urlTopicName.Id)
		unescaped, err := url.QueryUnescape(topicRaw)
		safeTopic := html_parser.EscapeWithEmojis(unescaped)
		if err != nil {
			r.Error(jerr.Get("error unescaping topic", err), http.StatusUnprocessableEntity)
			return
		}
		offset := r.Request.GetUrlParameterInt("offset")
		threads, err := db.GetThreads(uint(offset), unescaped)
		if err != nil {
			r.Error(jerr.Get("error getting threads from db", err), http.StatusInternalServerError)
			return
		}
		lastTopicList, err := cache.GetLastTopicList(r.Session.CookieId)
		if err != nil {
			jerr.Get("error getting last topic list", err).Print()
		}
		r.Helper["LastTopicList"] = lastTopicList
		r.Helper["Threads"] = threads
		r.Helper["Title"] = "Memo Topic - " + safeTopic
		r.Helper["Topic"] = safeTopic
		r.Helper["TopicEncoded"] = topicRaw
		res.SetPageAndOffset(r, offset)
		if len(threads) != 0 {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?", strings.TrimLeft(res.UrlTopicView + "/" + threads[0].GetUrlEncodedTopic() + res.UrlTopicThreads, "/"))
			r.Helper["TopicEncoded"] = threads[0].GetUrlEncodedTopic()
		}
		r.RenderTemplate(res.UrlTopicThreads)
	},
}
