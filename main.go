// Author Seth Hoenig 2015

// Command marathonctl provides total control over Marathon
// from the command line.
package main

import (
	"flag"
	"fmt"
	"os"
)

const Help = `marathonctl <flags...> [action] <args...>
 Actions
    app
       list                   - list all apps
       versions [id]          - list all versions of apps of id
       show [id] [version]    - show config of app of id and version
       create [jsonfile]      - deploy application defined in jsonfile
       update [id] [jsonfile] - update application id as defined in jsonfile
       restart [id]           - restart app of id
       destroy [id]           - destroy and remove all instances of id

    task
       list               - list all tasks
       list [id]          - list tasks of app of id
       kill [id]          - kill all tasks of app id
       kill [id] [taskid] - kill task taskid of app id
       queue              - list all queued tasks

    group
       group list              - list all groups
       group list [groupid]    - list apps in group of groupid
       group create [jsonfile] - create a group defined in jsonfile
       group update [jsonfile] - update group defined as defined in jsonfile
       group destroy [groupid] - destroy group of groupid

    deploy
       list               - list all active deploys
       destroy [deployid] - cancel deployment of [deployid]

    marathon
       leader   - get the current Marathon leader
       abdicate - force the current leader to relinquish control
       ping     - ping Marathon master host[s]

 Flags
  -c [config file]
  -h [host]
  -u [user:password] (separated by colon)
  -f [format]
       human  (simplified, default)
       json   (json on one line)
       jsonpp (json pretty printed)
`

func Usage() {
	fmt.Fprintln(os.Stderr, Help)
	os.Exit(1)
}

func main() {
	host, login, e := Config()

	if e != nil {
		Usage()
	}

	l := NewLogin(host, login)
	c := NewClient(l)
	app := &Category{
		actions: map[string]Action{
			"list":     AppList{c},
			"versions": AppVersions{c},
			"show":     AppShow{c},
			"create":   AppCreate{c},
			"update":   AppUpdate{c},
			"restart":  AppRestart{c},
			"destroy":  AppDestroy{c},
		},
	}
	task := &Category{
		actions: map[string]Action{
			"list":  TaskList{c},
			"kill":  TaskKill{c},
			"queue": TaskQueue{c},
		},
	}
	group := &Category{
		actions: map[string]Action{
			"list":    GroupList{c},
			"create":  GroupCreate{c},
			"destroy": GroupDestroy{c},
		},
	}
	deploy := &Category{
		actions: map[string]Action{
			"list":   DeployList{c},
			"cancel": DeployCancel{c},
		},
	}
	marathon := &Category{
		actions: map[string]Action{
			"leader":   Leader{c},
			"abdicate": Abdicate{c},
			"ping":     Ping{c},
		},
	}
	t := &Tool{
		selections: map[string]Selector{
			"app":      app,
			"task":     task,
			"group":    group,
			"deploy":   deploy,
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
