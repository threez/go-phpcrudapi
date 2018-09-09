package phpcrudapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	url string
	// basic auth
	username, password string
}

func New(url string) *Client {
	return &Client{url: url}
}

func (c *Client) BasicAuth(username, password string) {
	c.username = username
	c.password = password
}

func (c Client) All(ctx context.Context, v interface{}) error {
	return c.AllFiltered(ctx, &NoFilter{}, v)
}

func (c Client) AllFiltered(ctx context.Context, f Filter, v interface{}) error {
	st, table, err := getTypeAndTableSlice(v)
	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, "GET", "/"+table, f.Query(), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// parse json
	var inter map[string]tableResultSet
	err = json.NewDecoder(resp.Body).Decode(&inter)
	if err != nil {
		return err
	}

	return unmarshalSlice(st, table, inter, v)
}

func (c Client) newRequest(ctx context.Context, method, path, query string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	u.Path = path
	u.RawQuery = query

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}
