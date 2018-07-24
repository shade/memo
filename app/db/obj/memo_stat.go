package obj

import (
	"github.com/memocash/memo/app/obj/chart"
	"time"
)

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

func GetCohortStatData(memoCohortStats []MemoCohortStat, users bool) []chart.MultiSeries {
	var cohortActionsData = make(map[time.Time]map[time.Time]int)
	var cohorts []time.Time
	var dates []time.Time
	for _, cohortStat := range memoCohortStats {
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
		if users {
			cohortActionsData[cohortStat.Date][cohortStat.Cohort] = cohortStat.NumUsers
		} else {
			cohortActionsData[cohortStat.Date][cohortStat.Cohort] = cohortStat.NumPosts
		}
	}

	var multiSeriesMap = make(map[time.Time]chart.MultiSeries)
	for _, cohort := range cohorts {
		multiSeriesMap[cohort] = chart.MultiSeries{
			Name: cohort.Format("Jan 2006"),
		}
	}
	for _, date := range dates {
		cohortItems := cohortActionsData[date]
		for _, cohort := range cohorts {
			seriesCohort := multiSeriesMap[cohort]
			num, _ := cohortItems[cohort]
			seriesCohort.Data = append(seriesCohort.Data, []int64{date.Unix() * 1000, int64(num)})
			multiSeriesMap[cohort] = seriesCohort
		}
	}
	var multiSeries []chart.MultiSeries
	for _, cohort := range cohorts {
		multiSeries = append(multiSeries, multiSeriesMap[cohort])
	}

	return multiSeries
}
