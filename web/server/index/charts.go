package index

import (
	"encoding/json"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/db/view"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var chartsRoute = web.Route{
	Pattern: res.UrlCharts,
	Handler: func(r *web.Response) {
		userFirstPostStats, err := db.GetUserFirstPostStats()
		if err != nil {
			r.Error(jerr.Get("error getting first post stats", err), http.StatusInternalServerError)
			return
		}
		var userFirstPostData [][]int64
		for _, stat := range userFirstPostStats {
			userFirstPostData = append(userFirstPostData, []int64{
				stat.Date.Unix() * 1000,
				int64(stat.NumUsers),
			})
		}
		userFirstPostJson, err := json.Marshal(userFirstPostData)
		if err != nil {
			r.Error(jerr.Get("error marshalling user first posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["FirstPosts"] = string(userFirstPostJson)

		cohortStats, err := db.GetMemoCohortStats()
		if err != nil {
			r.Error(jerr.Get("error getting cohort stats", err), http.StatusInternalServerError)
			return
		}
		cohortActionsJson, err := json.Marshal(view.GetCohortStatData(cohortStats, false))
		if err != nil {
			r.Error(jerr.Get("error marshalling cohort actions", err), http.StatusInternalServerError)
			return
		}
		r.Helper["CohortActions"] = string(cohortActionsJson)

		cohortUsersJson, err := json.Marshal(view.GetCohortStatData(cohortStats, true))
		if err != nil {
			r.Error(jerr.Get("error marshalling cohort users", err), http.StatusInternalServerError)
			return
		}
		r.Helper["CohortUsers"] = string(cohortUsersJson)

		r.Helper["Title"] = "Memo - Charts"

		r.RenderTemplate(res.TmplCharts)
	},
}
