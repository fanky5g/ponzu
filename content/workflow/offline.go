package workflow

import "text/template"

type Offline struct{}

func (workflow *Offline) GetState() State {
	return OfflineState
}

func (workflow *Offline) GetValidTransitions() []Workflow {
	return []Workflow{
		&Offline{},
		&Archived{},
		&Online{},
	}
}

func (workflow *Offline) GetPastTransitions() []Workflow {
	return []Workflow{&Online{}}
}

func (workflow *Offline) GetAction(source Workflow) (*template.Template, error) {
	if source.GetState() == workflow.GetState() {
		return template.New("action").Parse("Re-unpublish {{ .EntityName }}")
	}

	return template.New("action").Parse("Unpublish {{ .EntityName }}")
}
