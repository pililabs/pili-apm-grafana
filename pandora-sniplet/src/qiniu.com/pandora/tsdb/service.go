package tsdb

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	. "qiniu.com/pandora/base"
	"qiniu.com/pandora/base/config"
	"qiniu.com/pandora/base/request"
)

var builder errBuilder

type Tsdb struct {
	Config     *config.Config
	HTTPClient *http.Client
}

func NewConfig() *config.Config {
	return config.NewConfig()
}

func New(c *config.Config) (TsdbAPI, error) {
	return newClient(c)
}

func newClient(c *config.Config) (p *Tsdb, err error) {
	if !strings.HasPrefix(c.Endpoint, "http://") && !strings.HasPrefix(c.Endpoint, "https://") {
		err = fmt.Errorf("endpoint should start with 'http://' or 'https://'")
		return
	}
	if strings.HasSuffix(c.Endpoint, "/") {
		err = fmt.Errorf("endpoint should not end with '/'")
		return
	}

	var t = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   c.DialTimeout,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: c.ResponseTimeout,
	}

	p = &Tsdb{
		Config:     c,
		HTTPClient: &http.Client{Transport: t},
	}

	return
}

func (c *Tsdb) newRequest(op *request.Operation, token string, v interface{}) *request.Request {
	req := request.New(c.Config, c.HTTPClient, op, token, builder, v)
	req.Data = v
	return req
}

func (c *Tsdb) newOperation(opName string, args ...interface{}) *request.Operation {
	var method, urlTmpl string
	switch opName {
	case OpCreateRepo:
		method, urlTmpl = MethodPost, "/v4/repos/%s"
	case OpListRepos:
		method, urlTmpl = MethodGet, "/v4/repos"
	case OpGetRepo:
		method, urlTmpl = MethodGet, "/v4/repos/%s"
	case OpDeleteRepo:
		method, urlTmpl = MethodDelete, "/v4/repos/%s"
	case OpUpdateRepoMetadata:
		method, urlTmpl = MethodPost, "/v4/repos/%s/meta"
	case OpDeleteRepoMetadata:
		method, urlTmpl = MethodDelete, "/v4/repos/%s/meta"
	case OpCreateSeries:
		method, urlTmpl = MethodPost, "/v4/repos/%s/series/%s"
	case OpUpdateSeriesMetadata:
		method, urlTmpl = MethodPost, "/v4/repos/%s/series/%s/meta"
	case OpDeleteSeriesMetadata:
		method, urlTmpl = MethodDelete, "/v4/repos/%s/series/%s/meta"
	case OpListSeries:
		method, urlTmpl = MethodGet, "/v4/repos/%s/series"
	case OpDeleteSeries:
		method, urlTmpl = MethodDelete, "/v4/repos/%s/series/%s"
	case OpCreateView:
		method, urlTmpl = MethodPost, "/v4/repos/%s/views/%s"
	case OpListView:
		method, urlTmpl = MethodGet, "/v4/repos/%s/views"
	case OpDeleteView:
		method, urlTmpl = MethodDelete, "/v4/repos/%s/views/%s"
	case OpGetView:
		method, urlTmpl = MethodGet, "/v4/repos/%s/views/%s"
	case OpQueryPoints:
		method, urlTmpl = MethodPost, "/v4/repos/%s/query"
	case OpWritePoints:
		method, urlTmpl = MethodPost, "/v4/repos/%s/points"
	default:
		c.Config.Logger.Errorf("unmatched operation name: %s", opName)
		return nil
	}

	return &request.Operation{
		Name:   opName,
		Method: method,
		Path:   fmt.Sprintf(urlTmpl, args...),
	}
}
