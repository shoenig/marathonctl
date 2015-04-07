package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// marathon [actions]

// ping (todo ping all hosts)
type Ping struct {
	client *Client
	format Formatter
}

func (p Ping) Apply(args []string) {
	request := p.client.GET("/ping")
	start := time.Now()
	response, e := p.client.Do(request)
	Check(e == nil, "ping failed", e)
	defer response.Body.Close()
	elapsed := fmt.Sprintf("%v", time.Now().Sub(start))
	fmt.Println(p.format.Format(strings.NewReader(elapsed), p.Humanize))
}

func (P Ping) Humanize(body io.Reader) string {
	b, e := ioutil.ReadAll(body)
	Check(e == nil, "reading ping response failed", e)
	// todo print HOST ELAPSED title and print hosts
	return fmt.Sprintf("elapsed: %s", string(b))
}

// leader
type Leader struct {
	client *Client
	format Formatter
}

func (l Leader) Apply(args []string) {
	request := l.client.GET("/v2/leader")
	response, e := l.client.Do(request)
	Check(e == nil, "get leader failed", e)
	c := response.StatusCode
	Check(c == 200, "get leader bad status", c)
	defer response.Body.Close()
	fmt.Println(l.format.Format(response.Body, l.Humanize))
}

func (l Leader) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var which Which
	e := dec.Decode(&which)
	Check(e == nil, "failed to decode response", e)
	text := "LEADER\n" + which.Leader
	return Columnize(text)
}

// abdicate
type Abdicate struct {
	client *Client
	format Formatter
}

func (a Abdicate) Apply(args []string) {
	request := a.client.DELETE("/v2/leader")
	response, e := a.client.Do(request)
	Check(e == nil, "abdicate request failed", e)
	c := response.StatusCode
	Check(c == 200, "abdicate bad status", c)
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a Abdicate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var mess Message
	e := dec.Decode(&mess)
	Check(e == nil, "failed to decode response", e)
	return "MESSAGE\n" + mess.Message
}
