package main

import (
	"fmt"
	"os"
)

func DoAction(m *Marathon, args []string) {
	if len(args) < 1 {
		fmt.Println("must specify an action")
		os.Exit(1)
	}

	action := args[0]

	switch action {
	case "create":
		fmt.Println("doing create")
	case "list":
		fmt.Println("doing list")
	default:
		fmt.Println("invalid action: " + action)
		os.Exit(1)
	}
}
