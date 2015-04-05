package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

// All actions under the task command

type TaskList struct {
	client *Client
}

func (t TaskList) Apply(args []string) {
	switch len(args) {
	case 0:
		t.listAll()
	case 1:
		t.listById(args[0])
	default:
		Check(false, "too many arguments")
	}
}

func (t TaskList) listAll() {
	path := "/v2/tasks"
	request := t.client.GET(path)
	response, e := t.client.Do(request)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var tasks Tasks
	e = dec.Decode(&tasks)
	Check(e == nil, "failed to unmarshal response", e)
	var b bytes.Buffer
	for _, task := range tasks.Tasks {
		b.WriteString(task.AppID)
		b.WriteString(" ")
		b.WriteString(task.Host)
		b.WriteString(" ")
		b.WriteString(task.Version)
		b.WriteString(" ")
		b.WriteString(task.ID)
		b.WriteString("\n")
	}
	title := "APPID HOST VERSION TASKID\n"
	text := title + b.String()
	fmt.Println(Columnize(text))
}

func (t TaskList) listById(id string) {
	esc := url.QueryEscape(id)
	path := "/v2/apps/" + esc + "?embed=apps.tasks"
	request := t.client.GET(path)
	response, e := t.client.Do(request)
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

type TaskKill struct {
	client *Client
}

func (t TaskKill) Apply(args []string) {

}
