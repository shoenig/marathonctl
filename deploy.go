package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// deploy [actions]

type DeployList struct {
	client *Client
}

func (d DeployList) Apply(args []string) {
	Check(len(args) == 0, "no arguments")
	request := d.client.GET("/v2/deployments")
	response, e := d.client.Do(request)
	Check(e == nil, "failed to get response")
	dec := json.NewDecoder(response.Body)
	defer response.Body.Close()
	var deploys Deploys
	e = dec.Decode(&deploys)
	Check(e == nil, "failed to unmarshal response", e)
	title := "DEPLOYID VERSION PROGRESS APPS\n"
	var b bytes.Buffer
	for _, deploy := range deploys {
		b.WriteString(deploy.DeployID)
		b.WriteString(" ")
		b.WriteString(deploy.Version)
		b.WriteString(" ")
		b.WriteString(strconv.Itoa(deploy.CurrentStep))
		b.WriteString("/")
		b.WriteString(strconv.Itoa(deploy.TotalSteps))
		b.WriteString(" ")
		for _, app := range deploy.AffectedApps {
			b.WriteString(app)
		}
		b.UnreadRune()
		b.WriteString("\n")
	}
	text := title + b.String()
	fmt.Println(Columnize(text))
}

type DeployCancel struct {
	client *Client
}

func (d DeployCancel) Apply(args []string) {
	Check(len(args) == 1, "must supply deployid")
	deployid := url.QueryEscape(args[0])
	path := "/v2/deployments/" + deployid
	request := d.client.DELETE(path)
	response, e := d.client.Do(request)
	Check(e == nil, "failed to cancel deploy", e)
	dec := json.NewDecoder(response.Body)
	var rollback Update
	e = dec.Decode(&rollback)
	Check(e == nil, "failed to decode response", e)
	title := "DEPLOYID VERSION\n"
	text := title + rollback.DeploymentID + " " + rollback.Version
	fmt.Println(Columnize(text))
}
