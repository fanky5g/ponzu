package content

import (
	"testing"

	"github.com/fanky5g/ponzu/internal/test/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ReferenceLoaderTestSuite struct {
	suite.Suite
}

func (s *ReferenceLoaderTestSuite) Test_buildReferences() {
	cases := []struct {
		Type                 interface{}
		ExpectedReferenceMap map[string][]string
	}{
		{
			Type: &types.Page{
				Title: "Home",
				URL:   "https://homepage.domain",
				ContentBlocks: &types.PageContentBlocks{
					{
						Type: "ImageGallery",
						Value: types.ImageAndTextBlock{
							Image: "5aa77561-d8f3-42b8-8d9c-98c44b7916ff",
						},
					},
				},
			},
			ExpectedReferenceMap: map[string][]string{
				"Upload": {"5aa77561-d8f3-42b8-8d9c-98c44b7916ff"},
			},
		},
		{
			Type: &types.Author{
				Name:  "Foo",
				Age:   50,
				Image: "4e034b54-9e6b-4eb4-bd41-51a4f3baf599",
			},
			ExpectedReferenceMap: map[string][]string{
				"Upload": {"4e034b54-9e6b-4eb4-bd41-51a4f3baf599"},
			},
		},
		{
			Type: &types.Story{
				Title: "A story about you",
				Body:  "Story Body",
				Author: types.Author{
					Name:  "Foo",
					Age:   50,
					Image: "4e034b54-9e6b-4eb4-bd41-51a4f3baf599",
				},
			},
			ExpectedReferenceMap: map[string][]string{
				"Upload": {"4e034b54-9e6b-4eb4-bd41-51a4f3baf599"},
			},
		},
		{
			Type: types.Page{
				Title: "Home",
				URL:   "https://homepage.domain",
				ContentBlocks: &types.PageContentBlocks{
					{
						Type: "ImageGallery",
						Value: types.ImageAndTextBlock{
							Image: "5aa77561-d8f3-42b8-8d9c-98c44b7916ff",
						},
					},
					{
						Type: "ImageGallery",
						Value: types.ImageAndTextBlock{
							Image: "03a24c40-c4cd-484e-9d71-8c26fb37c68e",
						},
					},
				},
			},
			ExpectedReferenceMap: map[string][]string{
				"Upload": {"5aa77561-d8f3-42b8-8d9c-98c44b7916ff", "03a24c40-c4cd-484e-9d71-8c26fb37c68e"},
			},
		},
		{
			Type: &types.Product{
				Name:        "Product Name",
				Description: "<p>Product Description</p>",
				Images: []string{
					"1ee2a499-12cc-4bd0-9394-5f918b7785cd",
					"35144bc4-e543-4423-9f18-e5d6f0d38176",
					"60cbe180-3610-47de-a1eb-ac2cdb6bf4c5",
				},
			},
			ExpectedReferenceMap: map[string][]string{
				"Photo": {
					"1ee2a499-12cc-4bd0-9394-5f918b7785cd",
					"35144bc4-e543-4423-9f18-e5d6f0d38176",
					"60cbe180-3610-47de-a1eb-ac2cdb6bf4c5",
				},
			},
		},
	}

	for i := range cases {
		assert.Equal(
			s.T(),
			cases[i].ExpectedReferenceMap,
			buildReferences(cases[i].Type),
		)
	}
}

func TestItem(t *testing.T) {
	suite.Run(t, new(ReferenceLoaderTestSuite))
}
