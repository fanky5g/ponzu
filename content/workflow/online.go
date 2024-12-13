package workflow

type Online struct{}

func (workflow *Online) GetState() State {
	return OnlineState
}

func (workflow *Online) GetValidTransitions() []Workflow {
	return []Workflow{
		&Online{},
		&Archived{},
		&Offline{},
	}
}

func (workflow *Online) GetPastTransitions() []Workflow {
	return []Workflow{&Draft{}}
}
