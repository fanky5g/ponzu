package workflow

import "text/template"

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

func (workflow *Preview) GetAction(source Workflow) (*template.Template, error) {
	if source.GetState() == workflow.GetState() {
		return template.New("action").Parse("Re-preview {{ .EntityName }}")
	}

	return template.New("action").Parse("Preview {{ .EntityName }}")
}
