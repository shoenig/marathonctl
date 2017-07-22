package main

// All actions under command marathon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ping (todo ping all hosts)
type MarathonPing struct {
	client *Client
	format Formatter
}

type MarathonPingResult struct {
	Host     string        `json:"host"`
	Duration time.Duration `json:"duration"`
}

type MarathonPingResultList []MarathonPingResult

func (p MarathonPing) Apply(args []string) {
	hosts := p.client.login.Hosts
	timings := make(map[string]time.Duration)
	for _, host := range hosts {
		request, e := http.NewRequest("GET", host+"/ping", nil)
		Check(e == nil, "could not create ping request", e)
		p.client.tweak(request)
		start := time.Now()
		_, err := p.client.client.Do(request)
		var elapsed time.Duration
		if err == nil {
			elapsed = time.Now().Sub(start)
		}
		timings[host] = elapsed
	}

	pingResultList := make([]MarathonPingResult, len(hosts))

	i := 0
	for host, duration := range timings {
		pingResultList[i] = MarathonPingResult{Host: host, Duration: duration}
		i += 1
	}
	jsonBytes, e := json.Marshal(pingResultList)
	Check(e == nil, "failed to marshal response", e)
	fmt.Println(p.format.Format(strings.NewReader(string(jsonBytes)), p.Humanize))
}

func (P MarathonPing) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var pingResults MarathonPingResultList
	e := dec.Decode(&pingResults)
	Check(e == nil, "failed to unmarshal response", e)
	title := "HOST DURATION\n"
	var b bytes.Buffer
	for _, pingResult := range pingResults {
		b.WriteString(pingResult.Host)
		b.WriteString(" ")
		b.WriteString(pingResult.Duration.String())
		b.WriteString("\n")
	}
	text := title + b.String()
	return Columnize(text)
}

// leader
type MarathonLeader struct {
	client *Client
	format Formatter
}

func (l MarathonLeader) Apply(args []string) {
	request := l.client.GET("/v2/leader")
	response, e := l.client.Do(request)
	Check(e == nil, "get leader failed", e)
	c := response.StatusCode
	Check(c == 200, "get leader bad status", c, l.format.Format(response.Body, l.Humanize))
	defer response.Body.Close()
	fmt.Println(l.format.Format(response.Body, l.Humanize))
}

func (l MarathonLeader) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var which Which
	e := dec.Decode(&which)
	Check(e == nil, "failed to decode response", e)
	text := "LEADER\n" + which.Leader
	return Columnize(text)
}

// abdicate
type MarathonAbdicate struct {
	client *Client
	format Formatter
}

func (a MarathonAbdicate) Apply(args []string) {
	request := a.client.DELETE("/v2/leader")
	response, e := a.client.Do(request)
	Check(e == nil, "abdicate request failed", e)
	c := response.StatusCode
	Check(c == 200, "abdicate bad status", c, a.format.Format(response.Body, a.Humanize))
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a MarathonAbdicate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var mess Message
	e := dec.Decode(&mess)
	Check(e == nil, "failed to decode response", e)
	return "MESSAGE\n" + mess.Message
}
