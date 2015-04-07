// Author Seth Hoenig 2015

// Command marathonctl is a CLI tool for Marathon
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
       group list                        - list all groups
       group list [groupid]              - list groups in groupid
       group create [jsonfile]           - create a group defined in jsonfile
       group update [groupid] [jsonfile] - update group groupid as defined in jsonfile
       group destroy [groupid]           - destroy group of groupid

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
       human  (simplified columns, default)
       json   (json on one line)
       jsonpp (json pretty printed)
       raw    (the exact response from Marathon)
`

func Usage() {
	fmt.Fprintln(os.Stderr, Help)
	os.Exit(1)
}

func main() {
	host, login, format, e := Config()

	if e != nil {
		Usage()
	}

	f := NewFormatter(format)
	l := NewLogin(host, login)
	c := NewClient(l)
	app := &Category{
		actions: map[string]Action{
			"list":     AppList{c, f},
			"versions": AppVersions{c, f},
			"show":     AppShow{c, f},
			"create":   AppCreate{c, f},
			"update":   AppUpdate{c, f},
			"restart":  AppRestart{c, f},
			"destroy":  AppDestroy{c, f},
		},
	}
	task := &Category{
		actions: map[string]Action{
			"list":  TaskList{c, f},
			"kill":  TaskKill{c, f},
			"queue": TaskQueue{c, f},
		},
	}
	group := &Category{
		actions: map[string]Action{
			"list":    GroupList{c, f},
			"create":  GroupCreate{c, f},
			"update":  GroupUpdate{c, f},
			"destroy": GroupDestroy{c, f},
		},
	}
	deploy := &Category{
		actions: map[string]Action{
			"list":   DeployList{c, f},
			"cancel": DeployCancel{c, f},
		},
	}
	marathon := &Category{
		actions: map[string]Action{
			"leader":   MarathonLeader{c, f},
			"abdicate": MarathonAbdicate{c, f},
			"ping":     MarathonPing{c, f},
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
