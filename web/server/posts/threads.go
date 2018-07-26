package posts

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	"net/http"
	"strings"
)

var threadsRoute = web.Route{
	Pattern: res.UrlPostsThreads,
	Handler: func(r *web.Response) {
		preHandler(r)
		offset := r.Request.GetUrlParameterInt("offset")
		threads, err := db.GetThreads(uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting threads from db", err), http.StatusInternalServerError)
			return
		}
		res.SetPageAndOffset(r, offset)
		r.Helper["OffsetLink"] = fmt.Sprintf("%s?", strings.TrimLeft(res.UrlPostsThreads, "/"))
		r.Helper["Threads"] = threads
		r.Helper["Title"] = "Memo - Threads"
		r.Render()
	},
}
