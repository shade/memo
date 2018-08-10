package metric

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func AddHttpRequest(url string, pattern string, requestTime time.Duration, code int) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	if c == nil {
		return nil
	}
	tags := []string{
		fmt.Sprintf("%s:%s", TagUrl, url),
		fmt.Sprintf("%s:%s", TagPattern, pattern),
		fmt.Sprintf("%s:%d", TagCode, code),
	}
	err = c.Incr(NameHttpRequest, tags, 1)
	if err != nil {
		return jerr.Get("error incrementing http request", err)
	}
	err = c.Gauge(NameHttpRequestTime, requestTime.Seconds(), tags, 1)
	if err != nil {
		return jerr.Get("error incrementing http request", err)
	}
	return nil
}
