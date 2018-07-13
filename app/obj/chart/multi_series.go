package chart

type MultiSeries struct {
	Name string `json:"name"`
	Data [][]int64  `json:"data"`
}
