package metric

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

func AddHttpRequest(url string, code int) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	tags := []string{
		fmt.Sprintf("%s:%s", TagUrl, url),
		fmt.Sprintf("%s:%d", TagCode, code),
	}
	err = c.Incr(HttpRequest, tags, 1)
	if err != nil {
		return jerr.Get("error incrementing http request", err)
	}
	return nil
}
