package main

import "fmt"

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
	fmt.Println("list", id, version)
}

// upload
type Upload struct{}

func (u Upload) Apply(args []string) {

}
