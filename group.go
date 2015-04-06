package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// group [actions]

type GroupList struct {
	client *Client
}

func (g GroupList) Apply(args []string) {
	request := g.client.GET("/v2/groups")
	response, e := g.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	// b, e := ioutil.ReadAll(response.Body)
	// Check(e == nil, "err", e)
	// fmt.Println(string(b))

	dec := json.NewDecoder(response.Body)
	var root Group
	e = dec.Decode(&root)
	Check(e == nil, "failed to unmarshal response", e)
	title := "GROUPID VERSION GROUPS APPS\n"
	var b bytes.Buffer
	gatherGroup(&root, &b)
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
	Check(false, "group create todo")
}

type GroupDestroy struct {
	client *Client
}

func (g GroupDestroy) Apply(args []string) {
	Check(false, "group destroy todo")
}
