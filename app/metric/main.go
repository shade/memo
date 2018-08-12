package metric

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/config"
)

const (
	NameHttpRequest         = "http_request"
	NameHttpRequestTime     = "http_request_time"
	NameMemoBroadcast       = "memo_broadcast"
	NameMemoSave            = "memo_save"
	NameTransactionSaveTime = "transaction_save_time"
	NamePostSearch          = "post_search"
)

const (
	TagUrl          = "url"
	TagPattern      = "pattern"
	TagResponseCode = "response_code"

	TagOutputType = "type"

	TagSearchTerm = "search_term"

	TagCmd    = "cmd"
	TagCode   = "code"
	TagReason = "reason"
)

var statsdClient *statsd.Client
var statsdDisabled bool

func getStatsd() (*statsd.Client, error) {
	if statsdDisabled {
		return nil, nil
	} else if statsdClient == nil {
		statsdConfig := config.GetStatsdConfig()
		if statsdConfig.Port == 0 || statsdConfig.Host == "" {
			statsdDisabled = true
			return nil, nil
		}
		var err error
		statsdClient, err = statsd.New(fmt.Sprintf("%s:%d", statsdConfig.Host, statsdConfig.Port))
		if err != nil {
			return nil, jerr.Get("error getting statsd client", err)
		}
		statsdClient.Namespace = fmt.Sprintf("%s.", statsdConfig.Namespace)
	}
	return statsdClient, nil
}
