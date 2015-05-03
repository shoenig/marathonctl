package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Action interface {
	// Apply runs the command with args and gets the json result.
	Apply(args []string)
}

type Login struct {
	Hosts []string
	User  string
	Pass  string
}

func (l *Login) NeedsAuth() bool {
	// Auth is only needed if User and Pass were set
	return l.User != "" && l.Pass != ""
}

func NewLogin(hosts, login string) *Login {
	hostlist := strings.Split(hosts, ",")
	var user, pass string
	if login != "" {
		toks := strings.SplitN(login, ":", 2)
		user = toks[0]
		pass = toks[1]
	}
	return &Login{
		Hosts: hostlist,
		User:  user,
		Pass:  pass,
	}
}

type Tool struct {
	selections map[string]Selector
}

func (t *Tool) Start(args []string) {
	if len(args) == 0 {
		Usage()
	}
	if selection, ok := t.selections[args[0]]; !ok {
		Usage()
	} else {
		selection.Select(args[1:])
	}
}

type Client struct {
	client http.Client
	login  *Login
}

func NewClient(login *Login) *Client {
	return &Client{
		client: http.Client{},
		login:  login,
	}
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	// this is not ghetto at all
	original := r.URL.String()
	// try each host until success or run out
	for _, host := range c.login.Hosts {
		fixed := strings.Replace(original, "HOST", host, 1)
		url, e := url.Parse(fixed)
		Check(e == nil, "could not parse fixed url", e)
		r.URL = url
		if response, e := c.client.Do(r); e == nil {
			return response, nil
		} else {
			fmt.Fprintf(os.Stderr, "request to %s failed\n", host)
			ourl, e := url.Parse(original)
			Check(e == nil, "could not parse original url")
			r.URL = ourl
		}
	}
	return nil, errors.New("requests to all hosts failed")
}

func (c *Client) GET(path string) *http.Request {
	url := "HOST" + path
	request, e := http.NewRequest("GET", url, nil)
	Check(e == nil, "failed to crete GET request", e)
	c.tweak(request)
	return request
}

func (c *Client) POST(path string, body io.ReadCloser) *http.Request {
	url := "HOST" + path
	request, e := http.NewRequest("POST", url, body)
	Check(e == nil, "failed to create POST request", e)
	c.tweak(request)
	return request
}

func (c *Client) DELETE(path string) *http.Request {
	url := "HOST" + path
	request, e := http.NewRequest("DELETE", url, nil)
	Check(e == nil, "failed to create DELETE request", e)
	c.tweak(request)
	return request
}

func (c *Client) PUT(path string, body io.ReadCloser) *http.Request {
	url := "HOST" + path
	request, e := http.NewRequest("PUT", url, body)
	Check(e == nil, "failed to create PUT request", e)
	c.tweak(request)
	return request
}

// tweak will set:
// Content-Type: application/json
// Accept: application/json
// Accept-Encoding: gzip, deflate, compress
func (c *Client) tweak(request *http.Request) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	if c.login.NeedsAuth() {
		request.SetBasicAuth(c.login.User, c.login.Pass)
	}
}
