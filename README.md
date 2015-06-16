# marathonctl
marathonctl is a command line tool for [Marathon](https://mesosphere.github.io/marathon/docs/rest-api.html)

## Install 
```
go get github.com/shoenig/marathonctl
```
- Maybe someday binary downloads will be available

Deployment via Docker
```
# Add to ~/.bash_aliases
alias marathonctl='docker run --rm --net=host shoenig/marathonctl:latest'
```

## Usage
```
marathonctl <flags...> [action] <args...>
 Actions
    app
       list                      - list all apps
       versions [id]             - list all versions of apps of id
       show [id]                 - show config and status of app of id (latest version)
       show [id] [version]       - show config and status of app of id and version
       create [jsonfile]         - deploy application defined in jsonfile
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

## Configuration
- Specify using "-c [file]" or "-h [host:port]"
    
### Configuration Properties
```
marathon.host: [hosts] (ex http://node1.indeed.com,http://node2.indeed.com,http://node3.indeed.com)
marathon.user: [user]
marathon.password: [password]
```

## Examples

#### Ping
- This example demonstrates -h, -u for host/login information
```
$ ./marathonctl -h http://marathon1:8080,http://marathon2:8080,http://marathon3:8080 -u user:pass marathon ping
HOST                   DURATION
http://marathon1:8080  11.004071ms
http://marathon2:8080  25.422ms
http://marathon3:8080  6.927772ms
```
#### Leader
- This example demonstrates -c and a file for host/login information
- This example demonstrates -f and the jsonpp (pretty printed json) output format
```
$ ./marathonctl -f jsonpp -c /etc/marathonctl.properties marathon leader
{
    "leader": "tst-mcontrol1:8080"
}
```
#### Abdicate
- This example demonstrates -f to specify one-line json output format
````
 ./marathonctl -f json -c /etc/marathonctl.properties marathon abdicate
{"message":"Leadership abdicted"}
````
#### Group List
- This example demonstrates the default human readable output
````
$ ./marathonctl -c /etc/marathonctl.properties group list
GROUPID                                     VERSION                   GROUPS  APPS  
/                                           2015-04-07T20:29:35.672Z  3       0     
/websites                                   2015-04-07T20:29:35.672Z  2       0     
/websites/indeed                            2015-04-07T20:29:35.672Z  1       0     
/websites/indeed/indeed-pings               2015-04-07T20:29:35.672Z  1       0     
/websites/indeed/indeed-pings/a.indeed.com  2015-04-07T20:29:35.672Z  0       2     
/websites/google                            2015-04-07T20:29:35.672Z  2       0     
/websites/google/news.google.com            2015-04-07T20:29:35.672Z  0       1     
/websites/google/calendar.google.com        2015-04-07T20:29:35.672Z  0       1
````
#### App Create
- This example demonstrates creating an app as specified in a json file
````
$ ./marathonctl -c /etc/marathonctl.properties app create sample/ping.google.json 
APPID                VERSION                   
/hoenig/ping-google  2015-04-07T21:41:53.440Z
````

## Bugs
- ping does not return json
