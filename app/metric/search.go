package metric

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

func AddMemoPostSearch(searchTerm string, pagePattern string) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	if c == nil {
		return nil
	}
	tags := []string{
		fmt.Sprintf("%s:%s", TagSearchTerm, searchTerm),
		fmt.Sprintf("%s:%s", TagPattern, pagePattern),
	}
	err = c.Incr(NamePostSearch, tags, 1)
	if err != nil {
		return jerr.Get("error incrementing post search", err)
	}
	return nil
}
