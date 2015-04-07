# marathonctl
marathonctl is a command line tool for [Marathon](https://mesosphere.github.io/marathon/docs/rest-api.html)

## Install 
```
go get github.com/shoenig/marathonctl
```
- Maybe someday binary downloads will be available

## Usage
```
marathonctl <flags...> [action] <args...>
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
       list                  - list all tasks
       list [id]             - list tasks of app of id
       kill [id]             - kill all tasks of app id
       kill [id] [taskid]    - kill task taskid of app id

    group
       group list              - list all groups
       group list [groupid]    - list apps in group of groupid
       group create [jsonfile] - create a group defined in jsonfile
       group update [jsonfile] - update group defined as defined in jsonfile
       group destroy [groupid] - destroy group of groupid

    deploy
       list               - list all active deploys
       queue              - list all queued deployes
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
```

## Bugs
- ping does not return json
