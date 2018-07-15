package index

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/feed_event"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var activityRoute = web.Route{
	Pattern: res.UrlActivity,
	Handler: func(r *web.Response) {
		setActivityFeed(r)
	},
}

func setActivityFeed(r *web.Response) {
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
	bal, err := cache.GetBalance(key.PkHash)
	if err != nil {
		r.Error(jerr.Get("error getting balance from cache", err), http.StatusInternalServerError)
		return
	}
	r.Helper["Balance"] = bal

	if bal == 0 {
		pf, err := profile.GetProfile(key.PkHash, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		err = pf.SetQr()
		if err != nil {
			r.Error(jerr.Get("error creating qr", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profile"] = pf
		r.RenderTemplate(res.TmplDashboardNoFunds)
		return
	}

	offset := r.Request.GetUrlParameterInt("offset")
	events, err := feed_event.GetEventsForUser(user.Id, key.PkHash, uint(offset))
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
	r.Helper["OffsetLink"] = fmt.Sprintf("%s?", res.UrlActivity)
	r.Helper["PrevOffset"] = prevOffset
	r.Helper["NextOffset"] = offset + 25

	r.RenderTemplate(res.TmplActivity)
}
