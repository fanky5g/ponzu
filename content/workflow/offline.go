package workflow

type Offline struct{}

func (workflow *Offline) GetState() State {
	return OfflineState
}

func (workflow *Offline) GetValidTransitions() []Workflow {
	return []Workflow{
		&Offline{},
		&Online{},
	}
}

func (workflow *Offline) GetPastTransitions() []Workflow {
	return []Workflow{&Online{}}
}
