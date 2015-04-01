package main

import (
	"encoding/json"
	"fmt"
)

func DoAction(c *Client, args []string) {
	Die(len(args) < 1, "must specify an action")

	action := args[0]

	switch action {
	case "create":
		Die(true, "create not supported yet")
	case "list":
		var id string
		var version string

		if len(args) > 1 {
			id = args[1]
		}

		if len(args) > 2 {
			version = args[2]
		}

		c.ListApps(id, version)

	default:
		Die(true, "unknown action", action)
	}
}

func (c *Client) ListApps(id, version string) {

	path := "/v2/apps"
	if id != "" {
		path += "/" + id
		if version != "" {
			path += "/versions/" + version
		}
	} // todo fix

	fmt.Println("path", path)

	request := c.GET(path)

	response, e := c.Do(request)
	Die(e != nil, "failed to get response", e)

	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var applications Applications
	e = dec.Decode(&applications)
	Die(e != nil, "failed to unmarshal response", e)

	applications.List()
}

// func CreateApp(m *Marathon, jsonfile string) {
// 	body, e := ioutil.ReadAll(jsonfile)
// 	Die(e != nil, "failed to read json file", e)

// 	c := http.Client{}
// 	url := m.Host + "/v2/apps" // do POST

// 	request, e := http.NewRequest("POST", url, nil)
// 	Die(e != nil, "failed to create request", e)
// }
