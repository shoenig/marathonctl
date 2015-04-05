// Author Seth Hoenig 2015

// Command marathonctl provides total control over Marathon
// from the command line.
package main

import (
	"flag"
	"fmt"
	"os"
)

const Howto = `marathonctl [-c conf] [-h host] [-u user:pass] <action>
 Actions
    list
       list                - lists instances
       list [id]           - lists instances of id
       list [id] [version] - lists instances of id and version

    versions
       versions [id] - list all versions of instances of id

    create
       create [jsonfile] - deploy application defined in jsonfile

    destroy
       destory [id] - destroy all instances of id

    group
       group list
       group list [groupid]
       group create [jsonfile]
       group update [jsonfile]
       group destroy [groupid]

    ping
       ping - ping any Marathon host

    leader
       leader - get the current Marathon leader

    abdicate
       abdicate - force the current leader to relinquish control
`

func main() {
	host, login, e := Config()

	if e != nil {
		Usage()
	}

	m := NewMarathon(host, login)
	c := NewClient(m)
	group := &Group{
		actions: map[string]Action{
			"list": Grouplist{c},
		},
	}
	t := &Tool{
		actions: map[string]Action{
			"list":     List{c},
			"create":   Create{c},
			"versions": Versions{c},
			"destroy":  Destroy{c},
			"group":    group,
			"ping":     Ping{c},
			"leader":   Leader{c},
			"abdicate": Abdicate{c},
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
