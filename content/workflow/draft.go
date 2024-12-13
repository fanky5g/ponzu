package workflow

type Draft struct{}

func (workflow *Draft) GetState() State {
	return DraftState
}

func (workflow *Draft) GetValidTransitions() []Workflow {
	return []Workflow{&Online{}}
}

func (workflow *Draft) GetPastTransitions() []Workflow {
	return nil
}
