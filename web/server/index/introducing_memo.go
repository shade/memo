package index

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var introducingMemoRoute = web.Route{
	Pattern: res.UrlIntroducing,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Introducing Memo"
		r.RenderTemplate(res.TmplIntroducing)
	},
}
