package generator

import "strings"

type Type int

const (
	Plain Type = iota + 1
	Content
	FieldCollection
)

type TypeDefinition struct {
	Name   string
	Label  string
	Blocks []Block
	Type
	Metadata
}

type Metadata struct {
	MethodReceiverName string
}

func NewTypeDefinition(itemType Type, args []string) (*TypeDefinition, error) {
	data := args[1:]
	blocks := make([]Block, len(data))
	for i, entry := range data {
		blockType := Field
		if itemType == FieldCollection {
			blockType = ContentBlock
		}

		blocks[i] = newBlock(entry, blockType)
	}

	name, label := parseName(args[0])

	return &TypeDefinition{
		Name:   name,
		Label:  label,
		Type:   itemType,
		Blocks: blocks,
		Metadata: Metadata{
			MethodReceiverName: strings.ToLower(string(name[0])),
		},
	}, nil
}
