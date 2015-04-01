package main

import (
	"strings"
)

type Marathon struct {
	host string // todo hosts
	user string
	pass string
}

func NewMarathon(host, login string) *Marathon {
	toks := strings.SplitN(login, ":", 2)
	return &Marathon{
		host: host,
		user: toks[0],
		pass: toks[1],
	}
}
