package main

type Selector interface {
	Select(args []string)
}

type Category struct {
	actions map[string]Action
}

func (c Category) Select(args []string) {
	Check(len(args) > 0, "must specify sub-action")
	if action, ok := c.actions[args[0]]; !ok {
		Usage()
	} else {
		action.Apply(args[1:])
	}
}
