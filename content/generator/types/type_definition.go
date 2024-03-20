package types

type TypeDefinition struct {
	Name          string
	Label         string
	Initial       string
	Fields        []Field
	ContentBlocks []ContentBlock
	HasReferences bool
}
