package metric

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/config"
)

const (
	NameHttpRequest   = "http_request"
	NameMemoBroadcast = "memo_broadcast"
	NameMemoSave = "memo_save"
)

var statsdClient *statsd.Client

func getStatsd() (*statsd.Client, error) {
	if statsdClient == nil {
		statsdConfig := config.GetStatsdConfig()
		var err error
		statsdClient, err = statsd.New("127.0.0.1:8125")
		if err != nil {
			return nil, jerr.Get("error getting statsd client", err)
		}
		statsdClient.Namespace = fmt.Sprintf("%s.", statsdConfig.Namespace)
	}
	return statsdClient, nil
}
