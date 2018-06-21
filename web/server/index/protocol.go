package index

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var protocolRoute = web.Route{
	Pattern: res.UrlProtocol,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Protocol"
		r.RenderTemplate(res.TmplProtocol)
	},
}
