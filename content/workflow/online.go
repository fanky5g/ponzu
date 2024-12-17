package workflow

import "text/template"

type Online struct{}

func (workflow *Online) GetState() State {
	return OnlineState
}

func (workflow *Online) GetValidTransitions() []Workflow {
	return []Workflow{
		&Online{},
		&Offline{},
	}
}

func (workflow *Online) GetPastTransitions() []Workflow {
	return []Workflow{&Draft{}}
}

func (workflow *Online) GetAction(source Workflow) (*template.Template, error) {
	if source.GetState() == workflow.GetState() {
		return template.New("action").Parse("Re-publish {{ .EntityName }}")
	}

	return template.New("action").Parse("Publish {{ .EntityName }}")
}
