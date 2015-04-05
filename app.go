package main

// app [actions]

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

type AppList struct {
	client *Client
}

func (a AppList) Apply(args []string) {
	var id string
	var version string

	if len(args) > 0 {
		id = args[0]
	} else if len(args) > 1 {
		version = args[1]
	}
	a.list(id, version)
}

func (a AppList) list(id, version string) {
	id = url.QueryEscape(id)

	path := "/v2/apps"

	if id != "" && version == "" {
		path += "/" + id + "?embed=apps.tasks" // why no work??
	} else if id != "" && version != "" {
		path += "/" + id + "/versions/" + version
	}

	fmt.Println("path", path)

	request := a.client.GET(path)

	response, e := a.client.Do(request)
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

type AppVersions struct {
	client *Client
}

func (a AppVersions) Apply(args []string) {
	Check(len(args) > 0, "must supply id")
	id := url.QueryEscape(args[0])
	path := "/v2/apps/" + id + "/versions"
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to list verions", e)
	defer response.Body.Close()
	b, e := ioutil.ReadAll(response.Body)
	Check(e == nil, "could not read response", e)
	ver := string(b)
	fmt.Println("versions = ", ver)
}

type AppShow struct {
	client *Client
}

func (s AppShow) Apply(args []string) {
	Check(false, "app show todo")
}

type AppCreate struct {
	client *Client
}

func (a AppCreate) Apply(args []string) {
	Check(len(args) == 1, "must specifiy 1 jsonfile")

	f, e := os.Open(args[0])
	Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()

	request := a.client.POST("/v2/apps", f)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var application Application
	e = dec.Decode(&application)
	Check(e == nil, "failed to decode response", e)
	fmt.Println(application.ID, application.Version)

}

type AppUpdate struct {
	client *Client
}

func (a AppUpdate) Apply(args []string) {
	Check(false, "app apply todo")
}

type AppRestart struct {
	client *Client
}

func (a AppRestart) Apply(args []string) {
	Check(false, "app restart todo")
}

type AppDestroy struct {
	client *Client
}

func (a AppDestroy) Apply(args []string) {
	Check(len(args) == 1, "must specify id")
	path := "/v2/apps/" + url.QueryEscape(args[0])
	request := a.client.DELETE(path)
	response, e := a.client.Do(request)
	Check(e == nil, "destroy app failed", e)
	c := response.StatusCode
	// documentation says this is 204, wtf
	Check(c == 200, "destroy app bad status", c)
	fmt.Println("destroyed", args[0])
}
