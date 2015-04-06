package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	switch len(args) {
	case 1:
		t.killAll(args[0])
	case 2:
		t.killOnly(args[0], args[1])
	default:
		Check(false, "task kill takes 1 or 2 arguments")
	}
}

func (t TaskKill) killAll(id string) {
	esc := url.QueryEscape(id)
	path := "/v2/apps/" + esc + "/tasks"
	request := t.client.DELETE(path)
	response, e := t.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()

	sc := response.StatusCode
	Check(sc != 404, "unknown id")
	Check(sc == 200, "failed with status code", sc)
	fmt.Println("success")
}

func (t TaskKill) killOnly(id, taskid string) {
	escID := url.QueryEscape(id)
	escTaskID := url.QueryEscape(taskid)
	path := "/v2/apps/" + escID + "/tasks/" + escTaskID
	request := t.client.DELETE(path)
	response, e := t.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	sc := response.StatusCode
	Check(sc != 404, "unknown appid or taskid")
	Check(sc == 200, "failed with status code", sc)
	fmt.Println("success")
}

type TaskQueue struct {
	client *Client
}

func (t TaskQueue) Apply(args []string) {
	Check(len(args) == 0, "no arguments")
	request := t.client.GET("/v2/queue")
	response, e := t.client.Do(request)
	Check(e == nil, "failed to get response", e)
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	var queue Queue
	e = dec.Decode(&queue)
	Check(e == nil, "failed to decode response", e)
	title := "APP VERSION OVERDUE\n"
	var b bytes.Buffer
	for _, queuedTask := range queue.Queue {
		b.WriteString(queuedTask.App.ID)
		b.WriteString(" ")
		b.WriteString(queuedTask.App.Version)
		b.WriteString(" ")
		b.WriteString(strconv.FormatBool(queuedTask.Delay["overdue"]))
		b.WriteString("\n")
	}
	text := title + b.String()
	fmt.Println(Columnize(text))
}
