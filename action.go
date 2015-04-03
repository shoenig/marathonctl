package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Action interface {
	Apply(args []string)
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
	id = url.QueryEscape(id) // todo whys no work??

	path := "/v2/apps"

	if id != "" && version == "" {
		path += "/" + id + "?embed=apps.tasks"
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

// upload
type Upload struct {
	client *Client
}

func (u Upload) Apply(args []string) {

}
