package index

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var needFundsRoute = web.Route{
	Pattern:    res.UrlNeedFunds,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
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
		r.Helper["Key"] = key
		r.RenderTemplate(res.TmplNeedFunds)
	},
}
