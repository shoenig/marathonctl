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
       list                      - list all apps
       versions [id]             - list all versions of apps of id
       show [id]                 - show config and status of app of id (latest version)
       show [id] [version]       - show config and status of app of id and version
       create [jsonfile]         - deploy application defined in jsonfile
       update [jsonfile]         - update application as defined in jsonfile
       update [id] [jsonfile]    - update application id as defined in jsonfile
       update cpu [id] [cpu%]    - update application id to have cpu% of cpu share
       update memory [id] [MB]   - update application id to have MB of memory
       update instances [id] [N] - update application id to have N instances
       restart [id]              - restart app of id
       destroy [id]              - destroy and remove all instances of id

    task
       list               - list all tasks
       list [id]          - list tasks of app of id
       kill [id]          - kill all tasks of app id
       kill [id] [taskid] - kill task taskid of app id
       queue              - list all queued tasks

    group
       list                        - list all groups
       list [groupid]              - list groups in groupid
       create [jsonfile]           - create a group defined in jsonfile
       update [groupid] [jsonfile] - update group groupid as defined in jsonfile
       destroy [groupid]           - destroy group of groupid

    deploy
       list               - list all active deploys
       destroy [deployid] - cancel deployment of [deployid]

    marathon
       leader   - get the current Marathon leader
       abdicate - force the current leader to relinquish control
       ping     - ping Marathon master host[s]

    artifact
       upload [path] [file]   - upload artifact to artifacts store
       get [path]             - get artifact from store
       delete [path]          - delete artifact from store

 Flags
  -c [config file]
  -h [host]
  -u [user:password] (separated by colon)
  -k - allow unverified tls connections
  -f [format]
       human  (simplified columns, default)
       json   (json on one line)
       jsonpp (json pretty printed)
       raw    (the exact response from Marathon)
  -v print git sha1
  -s print semver version
`

func Usage() {
	fmt.Fprintln(os.Stderr, Help)
	os.Exit(1)
}

func main() {
	conf, err := loadConfig()

	if err != nil {
		fmt.Printf("config error: %s\n\n", err)
		Usage()
	}

	if conf.version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if conf.semver {
		fmt.Println(Semver)
		os.Exit(0)
	}

	formatter := NewFormatter(conf.format)
	login := NewLogin(conf.host, conf.login)
	client := NewClient(login, conf.insecure)
	app := &Category{
		actions: map[string]Action{
			"list":     AppList{client, formatter},
			"versions": AppVersions{client, formatter},
			"show":     AppShow{client, formatter},
			"create":   AppCreate{client, formatter},
			"update":   AppUpdate{client, formatter},
			"restart":  AppRestart{client, formatter},
			"destroy":  AppDestroy{client, formatter},
		},
	}
	task := &Category{
		actions: map[string]Action{
			"list":  TaskList{client, formatter},
			"kill":  TaskKill{client, formatter},
			"queue": TaskQueue{client, formatter},
		},
	}
	group := &Category{
		actions: map[string]Action{
			"list":    GroupList{client, formatter},
			"create":  GroupCreate{client, formatter},
			"update":  GroupUpdate{client, formatter},
			"destroy": GroupDestroy{client, formatter},
		},
	}
	deploy := &Category{
		actions: map[string]Action{
			"list":    DeployList{client, formatter},
			"destroy": DeployCancel{client, formatter},
			"cancel":  DeployCancel{client, formatter},
		},
	}
	marathon := &Category{
		actions: map[string]Action{
			"leader":   MarathonLeader{client, formatter},
			"abdicate": MarathonAbdicate{client, formatter},
			"ping":     MarathonPing{client, formatter},
		},
	}
	artifact := &Category{
		actions: map[string]Action{
			"upload": ArtifactUpload{client, formatter},
			"get":    ArtifactGet{client, formatter},
			"delete": ArtifactDelete{client, formatter},
		},
	}
	t := &Tool{
		selections: map[string]Selector{
			"app":      app,
			"task":     task,
			"group":    group,
			"deploy":   deploy,
			"marathon": marathon,
			"artifact": artifact,
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
