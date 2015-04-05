package main

// app [actions]

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type AppList struct {
	client *Client
}

func (a AppList) Apply(args []string) {
	switch len(args) {
	case 0:
		a.listAll()
	case 1:
		a.listById(args[0])
	default:
		Check(false, "too many arguments")
	}
}

func (a AppList) listAll() {
	path := "/v2/apps"
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var applications Applications
	e = dec.Decode(&applications)
	Check(e == nil, "failed to unmarshal response", e)
	title := "APP VERSION USER\n"
	text := title + applications.String()
	fmt.Println(Columnize(text))
}

func (a AppList) listById(id string) {
	esc := url.QueryEscape(id)
	path := "/v2/apps/" + esc + "?embed=apps.tasks"
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var appbyid AppById
	e = dec.Decode(&appbyid)
	Check(e == nil, "failed to unmarshal response", e)

	var b bytes.Buffer
	for _, task := range appbyid.App.Tasks {
		b.WriteString(task.ID)
		b.WriteString(" ")
		b.WriteString(task.Host)
		b.WriteString(" ")
		b.WriteString(task.Version)
		// ports?
	}
	title := "ID HOST VERSION\n"
	text := title + b.String()
	fmt.Println(Columnize(text))
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
	dec := json.NewDecoder(response.Body)
	var versions Versions
	e = dec.Decode(&versions)
	Check(e == nil, "failed to unmarshal response", e)
	for _, version := range versions.Versions {
		fmt.Println(version)
	}
}

type AppShow struct {
	client *Client
}

func (a AppShow) Apply(args []string) {
	Check(len(args) == 2, "must provide id and version")
	id := url.QueryEscape(args[0])
	version := url.QueryEscape(args[1])
	path := "/v2/apps/" + id + "/versions/" + version
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to show app", e)
	dec := json.NewDecoder(response.Body)
	var application Application
	e = dec.Decode(&application)
	Check(e == nil, "failed to unmarshal response", e)
	title := "INSTANCES MEM CMD\n"
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(application.Instances))
	b.WriteString(" ")
	mem := fmt.Sprintf("%.2f", application.Mem)
	b.WriteString(mem)
	b.WriteString(" ")
	b.WriteString(application.Cmd)

	text := title + b.String()
	fmt.Println(Columnize(text))
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
	Check(response.StatusCode != 409, "app already exists")

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
	Check(len(args) == 2, "must specify id and jsonfile")
	id := url.QueryEscape(args[0])
	f, e := os.Open(args[1])
	Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()

	request := a.client.PUT("/v2/apps/"+id+"?force=true", f)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	sc := response.StatusCode
	Check(sc == 200, "bad status code", sc)

	dec := json.NewDecoder(response.Body)
	var update Update
	e = dec.Decode(&update)
	Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	fmt.Println(Columnize(text))
}

type AppRestart struct {
	client *Client
}

func (a AppRestart) Apply(args []string) {
	Check(len(args) == 1, "specify 1 app id to restart")
	id := url.QueryEscape(args[0])
	request := a.client.POST("/v2/apps/"+id+"/restart?force=true", nil)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var update Update
	e = dec.Decode(&update)
	Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	fmt.Println(Columnize(text))
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
