package main

import (
	"strings"
)

type Marathon struct {
	Host string // todo hosts
	User string
	Pass string
}

func NewMarathon(host, login string) *Marathon {
	toks := strings.SplitN(login, ":", 2)
	return &Marathon{
		Host: host,
		User: toks[0],
		Pass: toks[1],
	}
}
