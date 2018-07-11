package index

import (
	"encoding/json"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var chartsRoute = web.Route{
	Pattern: res.UrlCharts,
	Handler: func(r *web.Response) {
		stats, err := db.GetStats()
		if err != nil {
			r.Error(jerr.Get("error getting stats", err), http.StatusInternalServerError)
			return
		}
		var actionsData [][]int64
		for _, stat := range stats {
			actionsData = append(actionsData, []int64{
				stat.Date.Unix() * 1000,
				int64(stat.NumPosts),
			})
		}
		actionsJson, err := json.Marshal(actionsData)
		if err != nil {
			r.Error(jerr.Get("error marshalling actions", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Actions"] = string(actionsJson)
		var usersData [][]int64
		for _, stat := range stats {
			usersData = append(usersData, []int64{
				stat.Date.Unix() * 1000,
				int64(stat.NumUsers),
			})
		}
		usersJson, err := json.Marshal(usersData)
		if err != nil {
			r.Error(jerr.Get("error marshalling users", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Users"] = string(usersJson)
		r.RenderTemplate(res.TmplCharts)
	},
}
