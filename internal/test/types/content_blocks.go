package types

import (
	"encoding/json"
	"fmt"

	"github.com/fanky5g/ponzu/content"
)

type PageContentBlocks []content.FieldCollection

type ImageGallery struct {
	Headline   string `json:"headline"`
	Link       string `json:"link"`
	ButtonText string `json:"button_text"`
	Image      string `json:"image"`
}

type TextBlock struct {
	Text string `json:"text"`
}

type Link struct {
	ExternalUrl string `json:"external_url"`
	Label       string `json:"label"`
}

type ImageAndTextBlock struct {
	Image         string `json:"image" reference:"Upload"`
	ImagePosition string `json:"image_position"`
	Content       string `json:"entities"`
	Link          Link   `json:"link"`
}

func (p *PageContentBlocks) Name() string {
	return "Page Content Blocks"
}

func (p *PageContentBlocks) Data() []content.FieldCollection {
	return *p
}

func (p *PageContentBlocks) Add(fieldCollection content.FieldCollection) {
	*p = append(*p, fieldCollection)
}

func (p *PageContentBlocks) Set(i int, fieldCollection content.FieldCollection) {
	data := p.Data()
	data[i] = fieldCollection
	*p = data
}

func (p *PageContentBlocks) SetData(data []content.FieldCollection) {
	*p = data
}

func (p *PageContentBlocks) AllowedTypes() map[string]content.Builder {
	return map[string]content.Builder{
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

func (p *PageContentBlocks) UnmarshalJSON(b []byte) error {
	allowedTypes := p.AllowedTypes()

	var value []content.FieldCollection
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
