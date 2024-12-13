package workflow

type Preview struct{}

func (workflow *Preview) GetState() State {
	return PreviewState
}

func (workflow *Preview) GetValidTransitions() []Workflow {
	return []Workflow{
		&Preview{},
		&Online{},
	}
}

func (workflow *Preview) GetPastTransitions() []Workflow {
	return []Workflow{&Draft{}}
}
