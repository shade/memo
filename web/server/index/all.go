package index

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/feed_event"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var allRoute = web.Route{
	Pattern: res.UrlAll,
	Handler: func(r *web.Response) {
		r.Helper["Nav"] = "home"
		var userId uint
		var userPkHash []byte
		if auth.IsLoggedIn(r.Session.CookieId) {
			user, err := auth.GetSessionUser(r.Session.CookieId)
			if err != nil {
				r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
				return
			}
			key, err := db.GetKeyForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
				return
			}
			userId = user.Id
			userPkHash = key.PkHash
		}

		offset := r.Request.GetUrlParameterInt("offset")
		events, err := feed_event.GetAllEvents(userId, userPkHash, uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting events", err), http.StatusInternalServerError)
			return
		}
		r.Helper["FeedItems"] = events
		r.Helper["Offset"] = offset

		var prevOffset int
		if offset > 25 {
			prevOffset = offset - 25
		}
		page := offset/25 + 1
		r.Helper["Page"] = page
		r.Helper["OffsetLink"] = res.UrlAll
		r.Helper["PrevOffset"] = prevOffset
		r.Helper["NextOffset"] = offset + 25

		r.RenderTemplate(res.TmplAll)
	},
}
