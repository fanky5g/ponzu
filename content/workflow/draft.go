package workflow

type Draft struct{}

func (workflow *Draft) GetState() State {
	return DraftState
}

func (workflow *Draft) GetValidTransitions() []Workflow {
	return []Workflow{&Preview{}}
}

func (workflow *Draft) GetPastTransitions() []Workflow {
	return nil
}
