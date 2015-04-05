package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"
)

type Action interface {
	Apply(args []string)
}

// versions
type Versions struct {
	client *Client
}

func (v Versions) Apply(args []string) {
	Check(len(args) > 0, "must supply id")
	id := url.QueryEscape(args[0])
	path := "/v2/apps/" + id + "/versions"
	request := v.client.GET(path)
	response, e := v.client.Do(request)
	Check(e == nil, "failed to list verions", e)
	defer response.Body.Close()
	b, e := ioutil.ReadAll(response.Body)
	Check(e == nil, "could not read response", e)
	ver := string(b)
	fmt.Println("versions = ", ver)
}

// list
type List struct {
	client *Client
}

func (l List) Apply(args []string) {
	var id string
	var version string

	if len(args) > 0 {
		id = args[0]
	} else if len(args) > 1 {
		version = args[1]
	}
	l.list(id, version)
}

func (l List) list(id, version string) {
	id = url.QueryEscape(id)

	path := "/v2/apps"

	if id != "" && version == "" {
		path += "/" + id + "?embed=apps.tasks" // why no work??
	} else if id != "" && version != "" {
		path += "/" + id + "/versions/" + version
	}

	fmt.Println("path", path)

	request := l.client.GET(path)

	response, e := l.client.Do(request)
	Check(e == nil, "failed to get response", e)

	defer response.Body.Close()

	body, e := ioutil.ReadAll(response.Body)
	Check(e == nil, "failed to read body", e)
	fmt.Println(string(body))
	// dec := json.NewDecoder(response.Body)
	// var applications Applications
	// e = dec.Decode(&applications)
	// Check(e == nil, "failed to unmarshal response", e)
	// // fmt.Println(applications)

	// text := applications.String()
	// fmt.Println(Columnize(text))
}

// create
type Create struct {
	client *Client
}

func (u Create) Apply(args []string) {
	Check(len(args) == 1, "must specifiy 1 jsonfile")

	f, e := os.Open(args[0])
	Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()

	request := u.client.POST("/v2/apps", f)
	response, e := u.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var application Application
	e = dec.Decode(&application)
	Check(e == nil, "failed to decode response", e)
	fmt.Println(application.ID, application.Version)

}

// destroy
type Destroy struct {
	client *Client
}

func (d Destroy) Apply(args []string) {
	Check(len(args) == 1, "must specify id")
	path := "/v2/apps/" + url.QueryEscape(args[0])
	request := d.client.DELETE(path)
	response, e := d.client.Do(request)
	Check(e == nil, "destroy app failed", e)
	c := response.StatusCode
	// documentation says this is 204, wtf
	Check(c == 200, "destroy app bad status", c)
	fmt.Println("destroyed", args[0])
}

// ping
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

type Grouplist struct {
	client *Client
}

func (g Grouplist) Apply(args []string) {
	if len(args) == 0 {
		fmt.Println("list all the groups")
	} else {
		fmt.Println("list groups of id", args[0])
	}
}
