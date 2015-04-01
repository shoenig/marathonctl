// Author Seth Hoenig

// Usage
//     marathonctl [-c config] [-h host] [-u user:password] <action ...>
// Actions
//     list - lists all running applications
//     create <app config> - deploy application

// Command marathonctl gives you control over Marathon from the command line.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	host, login, e := Config()

	if e != nil {
		fmt.Println("failed to get configuration:", e)
		os.Exit(1)
	}

	fmt.Println(host, login)

	m := NewMarathon(host, login)
	DoAction(m, flag.Args())
}
