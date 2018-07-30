package metric

import (
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func AddTransactionSaveTime(duration time.Duration) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	err = c.Gauge(NameTransactionSaveTime, duration.Seconds(), nil, 1)
	if err != nil {
		return jerr.Get("error incrementing http request", err)
	}
	return nil
}
