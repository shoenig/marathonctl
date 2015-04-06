package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// group [actions]

type GroupList struct {
	client *Client
}

func (g GroupList) Apply(args []string) {
	switch len(args) {
	case 0:
		g.listGroups("")
	case 1:
		g.listGroups(args[0])
	default:
		Check(false, "expected 0 or 1 argument")
	}
}

func (g GroupList) listGroups(groupid string) {
	path := "/v2/groups"
	if groupid != "" {
		path += "/" + url.QueryEscape(groupid)
	}
	request := g.client.GET(path)
	response, e := g.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var root Group
	e = dec.Decode(&root)
	Check(e == nil, "failed to unmarshal response", e)
	printGroup(&root)
}

func printGroup(group *Group) {
	title := "GROUPID VERSION GROUPS APPS\n"
	var b bytes.Buffer
	gatherGroup(group, &b)
	text := title + b.String()
	fmt.Println(Columnize(text))
}

func gatherGroup(g *Group, b *bytes.Buffer) {
	b.WriteString(g.GroupID)
	b.WriteString(" ")
	b.WriteString(g.Version)
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(len(g.Groups)))
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(len(g.Apps)))
	b.WriteString("\n")
	for _, group := range g.Groups {
		gatherGroup(group, b)
	}
}

type GroupCreate struct {
	client *Client
}

func (g GroupCreate) Apply(args []string) {
	Check(len(args) == 1, "must supply 1 jsonfile")
	f, e := os.Open(args[0])
	Check(e == nil, "failed to open jsonfile", e)
	defer f.Close()
	request := g.client.POST("/v2/groups", f)
	response, e := g.client.Do(request)
	Check(e == nil, "failed to get response")
	defer response.Body.Close()
	Check(response.StatusCode != 409, "group already exists")

	dec := json.NewDecoder(response.Body)
	var update Update
	e = dec.Decode(&update)
	Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + update.DeploymentID + " " + update.Version
	fmt.Println(Columnize(text))
}

type GroupDestroy struct {
	client *Client
}

func (g GroupDestroy) Apply(args []string) {
	Check(len(args) == 1, "must specify groupid")
	groupid := url.QueryEscape(args[0])
	path := "/v2/groups/" + groupid
	request := g.client.DELETE(path)
	response, e := g.client.Do(request)
	Check(e == nil, "destroy group failed", e)
	defer response.Body.Close()
	c := response.StatusCode
	Check(c != 404, "unknown group")
	Check(c == 200, "destroy group bad status", c)
	dec := json.NewDecoder(response.Body)
	var versionmap map[string]string // ugh
	e = dec.Decode(&versionmap)
	Check(e == nil, "failed to decode response", e)

	v, ok := versionmap["version"]
	Check(ok, "version missing")

	fmt.Println("VERSION\n" + v)
}
