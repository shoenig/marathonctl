package main

// task [actions]

type TaskList struct {
	client *Client
}

func (t TaskList) Apply(args []string) {
	Check(false, "task list todo")
}

type TaskKill struct {
	client *Client
}

func (t TaskKill) Apply(args []string) {
	Check(false, "task apply todo")
}
