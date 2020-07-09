package datasource

import (
	"crypto/tls"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/fuloge/basework/configs"
	"net"
	"net/http"
	"time"
)

var (
	url     string
	logFile string
	esuser  string
	passwd  string
)

func init() {
	url = configs.EnvConfig.ES.Url
	logFile = configs.EnvConfig.ES.LogFile
	esuser = configs.EnvConfig.ES.User
	passwd = configs.EnvConfig.ES.Passwd
}

type ESClient struct {
	client *elastic.Client
}

func GetEsClient() (client *ESClient, err error) {
	esConfig := &elastic.Config{
		Addresses: []string{url},
		Username:  esuser,
		Password:  passwd,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}

	esclient, err := elastic.NewClient(*esConfig)
	if esclient != nil {
		return &ESClient{
			client: esclient,
		}, nil
	}

	return nil, err
}
