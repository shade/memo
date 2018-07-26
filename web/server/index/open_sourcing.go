package index

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var openSourcingMemoRoute = web.Route{
	Pattern: res.UrlOpenSource,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Open Sourcing Memo"
		r.RenderTemplate(res.TmplOpenSource)
	},
}
