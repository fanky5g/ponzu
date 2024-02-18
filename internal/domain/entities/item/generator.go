package item

type (
	TypeDefinition struct {
		Name          string
		Label         string
		Initial       string
		Fields        []Field
		HasReferences bool
	}

	Field struct {
		Name     string
		Label    string
		Initial  string
		TypeName string
		JSONName string
		ViewType string
		View     string

		IsReference       bool
		IsNested          bool
		ReferenceName     string
		ReferenceJSONTags []string
	}
)
