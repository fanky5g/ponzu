package workflow

type Archived struct{}

func (workflow *Archived) GetState() State {
	return ArchivedState
}

func (workflow *Archived) GetValidTransitions() []Workflow {
	return nil
}

func (workflow *Archived) GetPastTransitions() []Workflow {
	return []Workflow{&Offline{}}
}
