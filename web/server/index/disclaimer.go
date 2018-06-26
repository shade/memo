package index

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
)

var disclaimerRoute = web.Route{
	Pattern: res.UrlDisclaimer,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Disclaimer"
		r.RenderTemplate(res.TmplDisclaimer)
	},
}
