package contentgenerator

import (
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
	"github.com/fanky5g/ponzu/internal/domain/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go/format"
	"testing"
)

type GenerateTestSuite struct {
	suite.Suite
	gt *generator
}

func (s *GenerateTestSuite) SetupSuite() {
	var err error
	s.gt, err = setupGenerator()
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GenerateTestSuite) TestGenerateFieldCollection() {
	t := &item.TypeDefinition{
		Name:    "PageContentBlocks",
		Label:   "Page Content Blocks",
		Initial: "p",
		ContentBlocks: []item.ContentBlock{
			{
				TypeName: "ImageGallery",
				Label:    "Image Gallery",
			},
			{
				TypeName: "ImageAndTextBlock",
				Label:    "Image And Text Block",
			},
			{
				TypeName: "TextBlock",
				Label:    "Text Block",
			},
		},
	}

	expectedBuffer, err := format.Source([]byte(`
package content

import (
	"encoding/json"
	"fmt"
	"github.com/fanky5g/ponzu/internal/domain/entities/item"
)

type PageContentBlocks []item.FieldCollection

func (p *PageContentBlocks) Name() string {
	return "Page Content Blocks"
}

func (p *PageContentBlocks) Data() []item.FieldCollection {
	if p == nil {
		return nil
	}

	return *p
}

func (p *PageContentBlocks) AllowedTypes() map[string]item.EntityBuilder {
	return map[string]item.EntityBuilder{
		"ImageGallery": func() interface{} {
			return new(ImageGallery)
		},
		"ImageAndTextBlock": func() interface{} {
			return new(ImageAndTextBlock)
		},
		"TextBlock": func() interface{} {
			return new(TextBlock)
		},
	}
}

func (p *PageContentBlocks) Add(fieldCollection item.FieldCollection) {
	*p = append(*p, fieldCollection)
}

func (p *PageContentBlocks) Set(i int, fieldCollection item.FieldCollection) {
	data := p.Data()
	data[i] = fieldCollection
	*p = data
}

func (p *PageContentBlocks) SetData(data []item.FieldCollection) {
	*p = data
}

func (p *PageContentBlocks) UnmarshalJSON(b []byte) error {
	if p == nil {
		*p = make([]item.FieldCollection, 0)
	}

	allowedTypes := p.AllowedTypes()

	var value []item.FieldCollection
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	for i, t := range value {
		builder, ok := allowedTypes[t.Type]
		if !ok {
			return fmt.Errorf("type %s not implemented", t.Type)
		}

		entity := builder()
		byteRepresentation, err := json.Marshal(t.Value)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(byteRepresentation, entity); err != nil {
			return err
		}

		value[i].Value = entity
	}

	*p = value
	return nil
}

func init() {
	item.FieldCollectionTypes["PageContentBlocks"] = func() interface{} {
		return new(PageContentBlocks)
	}
}
`))

	if err != nil {
		s.T().Fatal(err)
	}

	buf, err := s.gt.generate(enum.TypeFieldCollection, t)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), expectedBuffer, buf)
	}
}

func TestGenerate(t *testing.T) {
	suite.Run(t, new(GenerateTestSuite))
}
