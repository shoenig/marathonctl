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
    app
       list                - lists apps
       list [id]           - lists apps of id
       list [id] [version] - lists apps of id and version
       versions [id] - list all versions of apps of id
       show [id] [version] - show config of app of id and version
       create [jsonfile] - deploy application defined in jsonfile
       update [jsonfile] - update application as defined in jsonfile
       restart [id] - restart app of id
       destroy [id] - destroy and remove all instances of id

    task
       list - list all tasks
       list [id] - list tasks of app of id
       kill [id] - kill all tasks of app id
       kill [id] [taskid] - kill task taskid of app id

    group
       group list
       group list [groupid]
       group create [jsonfile]
       group update [jsonfile]
       group destroy [groupid]

    deploy
       list - list all active deploys
       queue - list all queued deployes
       destroy [deployid] - cancel deployment of [deployid]

    marathon
       leader - get the current Marathon leader
       abdicate - force the current leader to relinquish control
       ping - ping Marathon host
`

func main() {
	host, login, e := Config()

	if e != nil {
		Usage()
	}

	l := NewLogin(host, login)
	c := NewClient(l)
	marathon := &Marathon{
		actions: map[string]Action{
			"leader":   Leader{c},
			"abdicate": Abdicate{c},
			"ping":     Ping{c},
		},
	}
	app := &App{
		actions: map[string]Action{
			"list":     List{c},
			"versions": Versions{c},
			"create":   Create{c},
			"destroy":  Destroy{c},
		},
	}
	group := &Group{
		actions: map[string]Action{
			"list": Grouplist{c},
		},
	}
	t := &Tool{
		actions: map[string]Action{
			"app":      app,
			"group":    group,
			"marathon": marathon,
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
