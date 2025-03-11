package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type TimeWebClient struct {
	host     string
	token    string
	basePath string
	Client   *http.Client
}

func NewTimeWebClient(basePath string, host string, token string) *TimeWebClient {
	return &TimeWebClient{
		host:     host,
		token:    token,
		basePath: basePath,
		Client:   &http.Client{},
	}
}

func (c *TimeWebClient) ServerInfo(ctx context.Context, serverID string) ([]byte, error) {
	return c.doRequest(ctx, fmt.Sprintf("/servers/%s", serverID), nil)
}

func (c *TimeWebClient) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't do request: %v", err)
		}
	}()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.URL.RawQuery = query.Encode()

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
