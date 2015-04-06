package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type DeployQueue struct {
	client *Client
}

func (d DeployQueue) Apply(args []string) {
	Check(false, "todo: deploy queue")
}

type DeployDestroy struct {
	client *Client
}

func (d DeployDestroy) Apply(args []string) {
	Check(false, "todo: deploy destroy")
}
