package main

// All actions under command app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type AppList struct {
	client *Client
	format Formatter
}

func (a AppList) Apply(args []string) {
	path := "/v2/apps"
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a AppList) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var applications Applications
	e := dec.Decode(&applications)
	Check(e == nil, "failed to unmarshal response", e)
	title := "APP VERSION USER\n"
	var b bytes.Buffer
	for _, app := range applications.Apps {
		b.WriteString(app.ID)
		b.WriteString(" ")
		b.WriteString(app.Version)
		b.WriteString(" ")
		b.WriteString(app.User)
		b.WriteString("\n")
	}
	text := title + b.String()
	return Columnize(text)
}

type AppVersions struct {
	client *Client
	format Formatter
}

func (a AppVersions) Apply(args []string) {
	Check(len(args) > 0, "must supply id")
	id := url.QueryEscape(args[0])
	path := "/v2/apps/" + id + "/versions"
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to list verions", e)
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a AppVersions) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var versions Versions
	e := dec.Decode(&versions)
	Check(e == nil, "failed to unmarshal response", e)
	var b bytes.Buffer
	b.WriteString("VERSIONS\n")
	for _, version := range versions.Versions {
		b.WriteString(version)
		b.WriteString("\n")
	}
	return b.String()
}

type AppShow struct {
	client *Client
	format Formatter
}

func (a AppShow) Apply(args []string) {
	path := ""
	fn := a.HumanizeById
	switch len(args) {
	case 1:
		id := url.QueryEscape(args[0])
		path = "/v2/apps/" + id
	case 2:
		id := url.QueryEscape(args[0])
		version := url.QueryEscape(args[1])
		path = "/v2/apps/" + id + "/versions/" + version
		fn = a.Humanize
	default:
		Check(false, "must provide id and/or version")
	}
	request := a.client.GET(path)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to show app", e)
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, fn))
}

func (a AppShow) HumanizeById(body io.Reader) string {
	dec := json.NewDecoder(body)
	var appbyid AppById
	e := dec.Decode(&appbyid)
	Check(e == nil, "failed to unmarshal response", e)
	title := "INSTANCES MEM CMD\n"
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(appbyid.App.Instances))
	b.WriteString(" ")
	mem := fmt.Sprintf("%.2f", appbyid.App.Mem)
	b.WriteString(mem)
	b.WriteString(" ")
	b.WriteString(appbyid.App.Cmd)
	text := title + b.String()
	return Columnize(text)
}

func (a AppShow) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var application Application
	e := dec.Decode(&application)
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
	return Columnize(text)
}

type AppCreate struct {
	client *Client
	format Formatter
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
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a AppCreate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var application Application
	e := dec.Decode(&application)
	Check(e == nil, "failed to decode response", e)
	var b bytes.Buffer
	b.WriteString("APPID VERSION\n")
	b.WriteString(application.ID)
	b.WriteString(" ")
	b.WriteString(application.Version)
	return Columnize(b.String())
}

type AppUpdate struct {
	client *Client
	format Formatter
}

func (a AppUpdate) Apply(args []string) {
	switch len(args) {
	case 2:
		a.fromJson(args)
	case 3:
		a.fromCLI(args)
	default:
		Check(false, "app update 2 or 3 arguments required")
	}
}

func (a AppUpdate) fromJson(args []string) {
	id := url.QueryEscape(args[0])
	f, e := os.Open(args[1])
	Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()
	a.update(id, f)
}

func (a AppUpdate) fromCLI(args []string) {
	val, e := strconv.ParseFloat(args[2], 64)
	Check(e == nil, "valid number required", args[1])
	var body string
	switch args[0] {
	case "instances":
		body = fmt.Sprintf("{\"instances\": %f}", val)
	case "memory":
		fallthrough
	case "mem":
		body = fmt.Sprintf("{\"mem\": %f}", val)
	case "cpu":
		body = fmt.Sprintf("{\"cpus\": %f}", val)
	default:
		Check(false, "unknown update option", args[0])
	}
	id := url.QueryEscape(args[1])
	a.update(id, ioutil.NopCloser(strings.NewReader(body)))
}

func (a AppUpdate) update(id string, body io.ReadCloser) {
	url := "/v2/apps/" + id + "?force=true"
	request := a.client.PUT(url, body)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	sc := response.StatusCode
	Check(sc == 200, "bad status code", sc)
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a AppUpdate) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var update Update
	e := dec.Decode(&update)
	Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	return Columnize(text)
}

type AppRestart struct {
	client *Client
	format Formatter
}

func (a AppRestart) Apply(args []string) {
	Check(len(args) == 1, "specify 1 app id to restart")
	id := url.QueryEscape(args[0])
	request := a.client.POST("/v2/apps/"+id+"/restart?force=true", nil)
	response, e := a.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a AppRestart) Humanize(body io.Reader) string {
	dec := json.NewDecoder(body)
	var update Update
	e := dec.Decode(&update)
	Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	return Columnize(text)
}

type AppDestroy struct {
	client *Client
	format Formatter
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
	defer response.Body.Close()
	fmt.Println(a.format.Format(response.Body, a.Humanize))
}

func (a AppDestroy) Humanize(body io.Reader) string {
	return "DESTROYED"
}
