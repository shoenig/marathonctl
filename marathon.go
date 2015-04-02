package main

import (
	"io"
	"net/http"
	"strings"
)

type Marathon struct {
	Host string // todo hosts
	User string
	Pass string
}

func NewMarathon(host, login string) *Marathon {
	toks := strings.SplitN(login, ":", 2)
	return &Marathon{
		Host: host,
		User: toks[0],
		Pass: toks[1],
	}
}

type Client struct {
	client   http.Client
	marathon *Marathon
}

func NewClient(m *Marathon) *Client {
	return &Client{
		client:   http.Client{},
		marathon: m,
	}
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	return c.client.Do(r)
}

func (c *Client) GET(path string) *http.Request {
	url := c.marathon.Host + path
	r, e := http.NewRequest("GET", url, nil)
	Die(e != nil, "failed to crete get request", e)
	r.SetBasicAuth(c.marathon.User, c.marathon.Pass)
	return r
}

func (c *Client) POST(path string, body io.ReadCloser) *http.Request {
	url := c.marathon.Host + path
	r, e := http.NewRequest("POST", url, body)
	Die(e != nil, "failed to create post request", e)
	r.Header.Set("Content-Type", "application/json")
	r.SetBasicAuth(c.marathon.User, c.marathon.Pass)
	return r
}
