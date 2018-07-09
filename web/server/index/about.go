package index

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var aboutRoute = web.Route{
	Pattern: res.UrlAbout,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - About"
		r.RenderTemplate(res.TmplAbout)
	},
}
