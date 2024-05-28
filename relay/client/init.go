package client

import (
	"net/http"
	"time"

	"github.com/songquanpeng/one-api/common/config"
)

var HTTPClient *http.Client
var ImpatientHTTPClient *http.Client

func init() {
	if config.RelayTimeout == 0 {
		HTTPClient = &http.Client{}
	} else {
		HTTPClient = &http.Client{
			Timeout: time.Duration(config.RelayTimeout) * time.Second,
		}
	}

	ImpatientHTTPClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}
