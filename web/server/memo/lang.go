package memo

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
	"net/http"
	"time"
	"github.com/memocash/memo/app/util"
	"github.com/jchavannes/jgo/jerr"
)

var setLangRoute = web.Route{
	Pattern:    res.UrlMemoSetLanguage,
	NeedsLogin: false,
	Handler: func(r *web.Response) {

		code := r.Request.GetFormValue("code")
		if !util.IsValidLang(code) {
			r.Error(jerr.New("unknown language"), http.StatusUnprocessableEntity)
			return
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "memo_language", Value: code, Path: "/", Expires: expiration, MaxAge: 31104000}
		http.SetCookie(r.Writer, &cookie)

		ref := r.Request.GetHeader("Referer")
		redir := "/"
		if len(ref) > 0 {
			redir = ref
		}
		r.SetRedirect(redir)
		r.Render()
	},
}

