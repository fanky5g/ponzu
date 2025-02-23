package generator

import (
	"regexp"
	"strings"
)

type BlockType int

const (
	Field BlockType = iota + 1
	ContentBlock
)

var (
	referenceRegex = regexp.MustCompile("^(?:\\[\\])?@")
)

type BlockDefinition struct {
	Title       string
	Type        string
	IsArray     bool
	IsReference bool
}

// Block is the building block(s) of types
type Block struct {
	Type          BlockType
	Name          string
	Label         string
	JSONName      string
	TypeName      string
	ReferenceName string
	Definition    BlockDefinition
}

func newBlock(definition string, kind BlockType) Block {
	data := strings.Split(definition, ":")
	title := strings.TrimSpace(data[0])
	blockType := strings.TrimSpace(strings.Join(data[1:], ":"))

	name, label := parseName(title)
	if kind == ContentBlock {
		name, label = parseName(definition)
		blockType = title
	}

	isReference := referenceRegex.MatchString(blockType)
	blockDefinition := BlockDefinition{
		Title:       title,
		Type:        blockType,
		IsArray:     strings.HasPrefix(blockType, "[]"),
		IsReference: referenceRegex.MatchString(blockType),
	}

	block := Block{
		Type:       kind,
		Name:       name,
		Label:      label,
		JSONName:   getJSONName(title),
		Definition: blockDefinition,
		TypeName:   getTypeName(blockDefinition),
	}

	if kind == ContentBlock {
		block.TypeName = block.Name
	}

	if isReference {
		block.ReferenceName = getReferenceName(blockDefinition)
	}

	return block
}

func getTypeName(blockDefinition BlockDefinition) string {
	if !blockDefinition.IsReference {
		return strings.ToLower(strings.Split(blockDefinition.Type, ":")[0])
	}

	if blockDefinition.IsArray {
		return "[]string"
	}

	return "string"
}

func getReferenceName(definition BlockDefinition) string {
	if !definition.IsReference {
		return ""
	}

	var referenceType string
	if definition.IsArray {
		referenceType = strings.TrimPrefix(definition.Type, "[]@")
	} else {
		referenceType = strings.TrimPrefix(definition.Type, "@")
	}

	referenceName, _ := parseName(referenceType)
	return referenceName
}
