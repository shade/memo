package obj

import (
	"fmt"
	"github.com/memocash/memo/app/obj/chart"
	"time"
)

type MemoStat struct {
	Date     time.Time
	NumPosts int
	NumUsers int
}

type MemoCohortStat struct {
	Date     time.Time
	Cohort   time.Time
	NumPosts int
	NumUsers int
}

func contains(item time.Time, items []time.Time) bool {
	for _, a := range items {
		if a.Equal(item) {
			return true
		}
	}
	return false
}

func GetCohortStatData(memoCohortStats []MemoCohortStat) []chart.MultiSeries {
	var cohortActionsData = make(map[time.Time]map[time.Time]int)
	var cohorts []time.Time
	var dates []time.Time
	for _, cohortStat := range memoCohortStats {
		fmt.Printf("date: %s\n", cohortStat.Date.String())
		_, ok := cohortActionsData[cohortStat.Date]
		if !ok {
			cohortActionsData[cohortStat.Date] = make(map[time.Time]int)
		}
		if ! contains(cohortStat.Cohort, cohorts) {
			cohorts = append(cohorts, cohortStat.Cohort)
		}
		if ! contains(cohortStat.Date, dates) {
			dates = append(dates, cohortStat.Date)
		}
		cohortActionsData[cohortStat.Date][cohortStat.Cohort] = cohortStat.NumPosts
	}

	var multiSeriesMap = make(map[time.Time]chart.MultiSeries)
	for _, cohort := range cohorts {
		multiSeriesMap[cohort] = chart.MultiSeries{
			Name: cohort.Format("Jan 2006"),
		}
	}
	for _, date := range dates {
		cohortPosts := cohortActionsData[date]
		for _, cohort := range cohorts {
			seriesCohort := multiSeriesMap[cohort]
			numPosts, _ := cohortPosts[cohort]
			seriesCohort.Data = append(seriesCohort.Data, []int64{date.Unix() * 1000, int64(numPosts)})
			multiSeriesMap[cohort] = seriesCohort
		}
	}
	var multiSeries []chart.MultiSeries
	for _, cohort := range cohorts {
		multiSeries = append(multiSeries, multiSeriesMap[cohort])
	}

	return multiSeries
}
