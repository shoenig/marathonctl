// Author Seth Hoenig 2015

// Usage
//     marathonctl [-c config] [-h host] [-u user:password] <action ...>
// Actions
//    -- list --
//      list                - lists all applications
//      list [id]           - lists instances of id
//      list [id] [version] - lists instances of id and version
//
//    -- versions --
//      versions [id] - list all running versions of id
//
//    -- create --
//      create [jsonfile] - deploy application defined in jsonfile
//
//    -- destroy --
//      destory [id] - destroy all instances of id

// Command marathonctl gives you control over Marathon from the command line.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	host, login, e := Config()

	Check(e == nil, "failed to get a Marathon configuration", e)

	m := NewMarathon(host, login)
	c := NewClient(m)
	t := &Tool{
		actions: map[string]Action{
			"list":     List{c},
			"create":   Create{c},
			"versions": Versions{c},
			"destroy":  Destroy{c},
			"ping":     Ping{c},
		},
	}

	t.Start(flag.Args())
}

func Check(b bool, args ...interface{}) {
	if !b {
		fmt.Fprintln(os.Stderr, args...)
		os.Exit(1)
	}
}
