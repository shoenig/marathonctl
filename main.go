// Author Seth Hoenig

// Usage
//     marathonctl [-c config] [-h host] [-u user:password] <action ...>
// Actions
//     -- list --
//     list            - lists all applications
//     list id         - lists applications of id
//     list id version - lists applications of id and version
//
//     -- create --
//     create <app json file> - deploy application

// Command marathonctl gives you control over Marathon from the command line.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	host, login, e := Config()

	Check(e == nil, "failed to get a mesos configuration", e)

	m := NewMarathon(host, login)
	c := NewClient(m)
	t := &Tool{
		client: c,
		actions: map[string]Action{
			"list":   List{c},
			"upload": Upload{c},
		},
	}

	t.Start(flag.Args())

	// DoAction(c, flag.Args())
}

func Check(b bool, args ...interface{}) {
	if !b {
		fmt.Fprintln(os.Stderr, args...)
		os.Exit(1)
	}
}
