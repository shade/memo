package index

import (
	"encoding/json"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/db/obj"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var chartsRoute = web.Route{
	Pattern: res.UrlCharts,
	Handler: func(r *web.Response) {
		stats, err := db.GetMemoStats()
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

		userLastPostStats, err := db.GetUserLastPostStats()
		if err != nil {
			r.Error(jerr.Get("error getting last post stats", err), http.StatusInternalServerError)
			return
		}
		var userLastPostData [][]int64
		for _, stat := range userLastPostStats {
			userLastPostData = append(userLastPostData, []int64{
				stat.Date.Unix() * 1000,
				int64(stat.NumUsers),
			})
		}
		userLastPostJson, err := json.Marshal(userLastPostData)
		if err != nil {
			r.Error(jerr.Get("error marshalling user last posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["LastPosts"] = string(userLastPostJson)

		cohortStats, err := db.GetMemoCohortStats()
		if err != nil {
			r.Error(jerr.Get("error getting cohort stats", err), http.StatusInternalServerError)
			return
		}
		cohortActionsData := obj.GetCohortStatData(cohortStats)
		cohortActionsJson, err := json.Marshal(cohortActionsData)
		if err != nil {
			r.Error(jerr.Get("error marshalling cohort actions", err), http.StatusInternalServerError)
			return
		}
		r.Helper["CohortActions"] = string(cohortActionsJson)

		r.RenderTemplate(res.TmplCharts)
	},
}
