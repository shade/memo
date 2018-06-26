package index

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var guidesRoute = web.Route{
	Pattern: res.UrlGuides,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Guides"
		r.RenderTemplate(res.TmplGuides)
	},
}
