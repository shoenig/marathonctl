package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
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

	dec := json.NewDecoder(response.Body)
	var applications Applications
	e = dec.Decode(&applications)
	Check(e == nil, "failed to unmarshal response", e)
	// fmt.Println(applications)

	text := applications.String()
	fmt.Println(Columnize(text))
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
