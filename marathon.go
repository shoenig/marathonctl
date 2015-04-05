package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// marathon [actions]

// ping (todo ping all hosts)
type Ping struct {
	client *Client
}

func (p Ping) Apply(args []string) {
	request := p.client.GET("/ping")
	start := time.Now()
	response, e := p.client.Do(request)
	elapsed := time.Now().Sub(start)
	Check(e == nil, "ping failed", e)
	c := response.StatusCode
	Check(c == 200, "ping bad status", c)
	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	Check(e == nil, "reading ping response failed", e)
	pong := strings.TrimSpace(string(body))
	fmt.Println(pong, elapsed)
}

// leader
type Leader struct {
	client *Client
}

func (l Leader) Apply(args []string) {
	request := l.client.GET("/v2/leader")
	response, e := l.client.Do(request)
	Check(e == nil, "get leader failed", e)
	c := response.StatusCode
	Check(c == 200, "get leader bad status", c)
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var which Which
	e = dec.Decode(&which)
	Check(e == nil, "failed to decode response", e)
	fmt.Println(which.Leader)
}

// abdicate
type Abdicate struct {
	client *Client
}

func (a Abdicate) Apply(args []string) {
	request := a.client.DELETE("/v2/leader")
	response, e := a.client.Do(request)
	Check(e == nil, "abdicate request failed", e)
	c := response.StatusCode
	Check(c == 200, "abdicate bad status", c)
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var mess Message
	e = dec.Decode(&mess)
	Check(e == nil, "failed to decode response", e)
	fmt.Println(mess.Message)
}
