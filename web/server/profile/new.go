package profile

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var newRoute = web.Route{
	Pattern:    res.UrlProfilesNew,
	Handler: func(r *web.Response) {
		profilesByDate(r, false)
		r.RenderTemplate(res.TmplProfilesNew)
	},
}
