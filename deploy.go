package main

// deploy [actions]

type DeployList struct {
	client *Client
}

func (d DeployList) Apply(args []string) {
	Check(false, "todo: deploy list")
}

type DeployQueue struct {
	client *Client
}

func (d DeployQueue) Apply(args []string) {
	Check(false, "todo: deploy queue")
}

type DeployDestroy struct {
	client *Client
}

func (d DeployDestroy) Apply(args []string) {
	Check(false, "todo: deploy destroy")
}
