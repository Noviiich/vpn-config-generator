package api

import (
	"context"
	"net/http"
	"net/url"
	"path"
)

type TimeWebClient struct {
	host     string
	basePath string
	Client   *http.Client
}

func NewTimeWebClient(basePath string, host string) *TimeWebClient {
	return &TimeWebClient{
		host:     host,
		basePath: basePath,
		Client:   &http.Client{},
	}
}

func (c *TimeWebClient) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
}
