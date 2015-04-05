package main

// group [actions]

type GroupList struct {
	client *Client
}

func (g GroupList) Apply(args []string) {
	Check(false, "group list todo")
}

type GroupCreate struct {
	client *Client
}

func (g GroupCreate) Apply(args []string) {
	Check(false, "group create todo")
}

type GroupDestroy struct {
	client *Client
}

func (g GroupDestroy) Apply(args []string) {
	Check(false, "group destroy todo")
}
