package content

import "github.com/fanky5g/ponzu/internal/search"

type Query struct {
	ID   string
	Type string
}

type TransitionInput struct {
	Query
	TargetState string
}

type Reference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListReferencesOutputResource struct {
	References []interface{} `json:"references"`
	Size       int           `json:"size"`
}

type GetReferenceInputResource struct {
	Type string
	ID   string
}

type ListReferencesInputResource struct {
	Type   string
	Search *search.Search
}
