package workflow

import "text/template"

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

func (workflow *Archived) GetAction(source Workflow) (*template.Template, error) {
	if source.GetState() == workflow.GetState() {
		return template.New("action").Parse("Re-archive {{ .EntityName }}")
	}

	return template.New("action").Parse("Archive {{ .EntityName }}")
}
