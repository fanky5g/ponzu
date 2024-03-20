package types

type Field struct {
	Name     string
	Label    string
	Initial  string
	TypeName string
	JSONName string
	ViewType string
	View     string

	IsReference       bool
	IsNested          bool
	IsFieldCollection bool
	ReferenceName     string
	ReferenceJSONTags []string
}
