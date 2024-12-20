package generator

type Generator interface {
	Generate(definition *TypeDefinition, writer Writer) error
	Initialize(definition *TypeDefinition, writer Writer) error
}
