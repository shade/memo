package metric

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

func AddMemoReject(cmd string, code string, reason string) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	if c == nil {
		return nil
	}
	tags := []string{
		fmt.Sprintf("%s:%s", TagCmd, cmd),
		fmt.Sprintf("%s:%s", TagCode, code),
		fmt.Sprintf("%s:%s", TagReason, reason),
	}
	err = c.Incr(NamePostSearch, tags, 1)
	if err != nil {
		return jerr.Get("error incrementing memo reject", err)
	}
	return nil
}
