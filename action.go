package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func DoAction(m *Marathon, args []string) {
	if len(args) < 1 {
		fmt.Println("must specify an action")
		os.Exit(1)
	}

	action := args[0]

	switch action {
	case "create":
		fmt.Println("doing create")
	case "list":
		fmt.Println("doing list")
		ListApps(m)
	default:
		fmt.Println("invalid action: " + action)
		os.Exit(1)
	}
}

func ListApps(m *Marathon) {
	c := http.Client{}
	request, e := http.NewRequest("GET", m.Host+"/v2/apps", nil)
	if e != nil {
		fmt.Println("failed to create request:", e)
		os.Exit(1)
	}
	request.SetBasicAuth(m.User, m.Pass)
	response, e := c.Do(request)
	if e != nil {
		fmt.Println("failed to get response:", e)
		os.Exit(1)
	}

	defer response.Body.Close()
	body, e := ioutil.ReadAll(response.Body)
	if e != nil {
		fmt.Println("failed to read response:", e)
		os.Exit(1)
	}

	fmt.Println(string(body))
}
