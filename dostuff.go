package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

func DoCreate(c *Client, args []string) {
	Check(len(args) < 2, "must supply jsonfile")
	fmt.Println("create", args[1])
	c.CreateApp(args[1])
}

func DoList(c *Client, args []string) {
	var id string
	var version string

	if len(args) > 1 {
		id = args[1]
	}

	if len(args) > 2 {
		version = args[2]
	}

	c.ListApps(id, version)
}

func DoVersions(c *Client, args []string) {
	Check(len(args) < 2, "must specifiy id")
	id := args[1]
	c.ListVersions(id)
}

func (c *Client) ListVersions(id string) {
	path := "/v2/apps/" + id + "/versions"
	request := c.GET(path)
	response, e := c.Do(request)
	Check(e == nil, "failed to list verions", e)
	defer response.Body.Close()
	b, e := ioutil.ReadAll(response.Body)
	Check(e == nil, "could not read response", e)
	v := string(b)
	fmt.Println("versions = ", v)
}

func (c *Client) ListApps(id, version string) {
	id = url.QueryEscape(id) // todo whys no work??

	path := "/v2/apps"

	if id != "" && version == "" {
		path += "/" + id + "?embed=apps.tasks"
	} else if id != "" && version != "" {
		path += "/" + id + "/versions/" + version
	}

	fmt.Println("path", path)

	request := c.GET(path)

	response, e := c.Do(request)
	Check(e == nil, "failed to get response", e)

	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var applications Applications
	e = dec.Decode(&applications)
	Check(e == nil, "failed to unmarshal response", e)
	fmt.Println(applications)

	applications.List()
}

func (c *Client) CreateApp(jsonfile string) {
	f, e := os.Open(jsonfile)
	Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()

	request := c.POST("/v2/apps", f)
	response, e := c.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	b, e := ioutil.ReadAll(response.Body)
	Check(e == nil, "failed", e)

	fmt.Println(string(b))

	// dec := json.NewDecoder(response.Body)
	// var application Application
	// e = dec.Decode(&application)
	// Die(e != nil, "failed to decode response", e)
	// fmt.Println(application)
}
