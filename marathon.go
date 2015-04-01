package main

import (
	"fmt"
	"net/http"
	"os"
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
	if e != nil {
		fmt.Println("creation of get request failed:", e)
		os.Exit(1)
	}
	r.SetBasicAuth(c.marathon.User, c.marathon.Pass)
	return r
}

// func (c *Client) POST() *http.Request {
// 	return nil
// }
