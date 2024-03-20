package content

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FieldCollectionTestSuite struct {
	suite.Suite
}

type pageContentBlocks []FieldCollection

type page struct {
	Title         string            `json:"title"`
	URL           string            `json:"url"`
	ContentBlocks pageContentBlocks `json:"content_blocks"`
}

type imageGallery struct {
	Headline   string `json:"headline"`
	Link       string `json:"link"`
	ButtonText string `json:"button_text"`
	Image      string `json:"image"`
}

type textBlock struct {
	Text string `json:"text"`
}

type link struct {
	ExternalUrl string `json:"external_url"`
	Label       string `json:"label"`
}

type imageAndTextBlock struct {
	Image         string `json:"image"`
	ImagePosition string `json:"image_position"`
	Content       string `json:"entities"`
	Link          link   `json:"link"`
}

func (p *pageContentBlocks) Name() string {
	return "Page Content Blocks"
}

func (p *pageContentBlocks) Data() []FieldCollection {
	return *p
}

func (p *pageContentBlocks) AllowedTypes() map[string]Builder {
	return map[string]Builder{
		"ImageGallery": func() interface{} {
			return new(imageGallery)
		},
		"ImageAndTextBlock": func() interface{} {
			return new(imageAndTextBlock)
		},
		"TextBlock": func() interface{} {
			return new(textBlock)
		},
	}
}

func (p *pageContentBlocks) UnmarshalJSON(b []byte) error {
	allowedTypes := p.AllowedTypes()

	var value []FieldCollection
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

func (s *FieldCollectionTestSuite) TestMarshal() {
	p := &page{
		Title: "Home",
		URL:   "https://ponzu.domain",
		ContentBlocks: []FieldCollection{
			{
				Type: "TextBlock",
				Value: &textBlock{
					Text: "This is some random text",
				},
			},
			{
				Type: "ImageAndTextBlock",
				Value: &imageAndTextBlock{
					Image:         "https://dummyimage.com/600x400/000/fff",
					ImagePosition: "right",
					Content:       "Feel the spirit",
					Link: link{
						ExternalUrl: "https://www.w3.org/Provider/Style/dummy.html",
					},
				},
			},
			{
				Type: "ImageGallery",
				Value: &imageGallery{
					Headline:   "Headline",
					Link:       "https://www.w3.org/Provider/Style/dummy.html",
					ButtonText: "Button",
					Image:      "https://dummyimage.com/600x400/000/fff",
				},
			},
		},
	}

	jsonValue, err := json.Marshal(p)
	if assert.NoError(s.T(), err) {
		assert.Equal(s.T(), []byte(`{"title":"Home","url":"https://ponzu.domain","content_blocks":[{"type":"TextBlock","value":{"text":"This is some random text"}},{"type":"ImageAndTextBlock","value":{"image":"https://dummyimage.com/600x400/000/fff","image_position":"right","entities":"Feel the spirit","link":{"external_url":"https://www.w3.org/Provider/Style/dummy.html","label":""}}},{"type":"ImageGallery","value":{"headline":"Headline","link":"https://www.w3.org/Provider/Style/dummy.html","button_text":"Button","image":"https://dummyimage.com/600x400/000/fff"}}]}`), jsonValue)
	}
}

func (s *FieldCollectionTestSuite) TestUnmarshalJSON() {
	data := "{\"title\":\"Home\",\"url\":\"https://ponzu.domain\",\"content_blocks\":[{\"type\":\"TextBlock\",\"value\":{\"text\":\"This is some random text\"}},{\"type\":\"ImageAndTextBlock\",\"value\":{\"image\":\"https://dummyimage.com/600x400/000/fff\",\"image_position\":\"right\",\"entities\":\"Feel the spirit\",\"link\":{\"external_url\":\"https://www.w3.org/Provider/Style/dummy.html\",\"label\":\"\"}}},{\"type\":\"ImageGallery\",\"value\":{\"headline\":\"Headline\",\"link\":\"https://www.w3.org/Provider/Style/dummy.html\",\"button_text\":\"Button\",\"image\":\"https://dummyimage.com/600x400/000/fff\"}}]}"
	expectedPage := page{
		Title: "Home",
		URL:   "https://ponzu.domain",
		ContentBlocks: []FieldCollection{
			{
				Type: "TextBlock",
				Value: &textBlock{
					Text: "This is some random text",
				},
			},
			{
				Type: "ImageAndTextBlock",
				Value: &imageAndTextBlock{
					Image:         "https://dummyimage.com/600x400/000/fff",
					ImagePosition: "right",
					Content:       "Feel the spirit",
					Link: link{
						ExternalUrl: "https://www.w3.org/Provider/Style/dummy.html",
					},
				},
			},
			{
				Type: "ImageGallery",
				Value: &imageGallery{
					Headline:   "Headline",
					Link:       "https://www.w3.org/Provider/Style/dummy.html",
					ButtonText: "Button",
					Image:      "https://dummyimage.com/600x400/000/fff",
				},
			},
		},
	}

	var p page
	if assert.NoError(s.T(), json.Unmarshal([]byte(data), &p)) {
		assert.Equal(s.T(), expectedPage, p)
	}
}

func TestFieldCollection(t *testing.T) {
	suite.Run(t, new(FieldCollectionTestSuite))
}
