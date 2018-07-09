package metric

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
)

const (
	TagSpendOutputType = "type"
)

func AddMemoBroadcast(outputType memo.OutputType) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	tags := []string{
		fmt.Sprintf("%s:%s", TagSpendOutputType, outputType.String()),
	}
	err = c.Incr(NameMemoBroadcast, tags, 1)
	if err != nil {
		return jerr.Get("error incrementing memo broadcast", err)
	}
	return nil
}

func AddMemoSave(code byte) error {
	c, err := getStatsd()
	if err != nil {
		return jerr.Get("error getting statsd", err)
	}
	tags := []string{
		fmt.Sprintf("%s:%s", TagSpendOutputType, memo.GetCodeString(code)),
	}
	err = c.Incr(NameMemoSave, tags, 1)
	if err != nil {
		return jerr.Get("error incrementing memo broadcast", err)
	}
	return nil
}
