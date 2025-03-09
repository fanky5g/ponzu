package editor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakePositionalPlaceholder(t *testing.T) {
	tt := []struct {
		name                          string
		tagName                       string
		expectedPositionalPlaceholder string
	}{
		{
			name:                          "Regular FieldNames",
			tagName:                       "author",
			expectedPositionalPlaceholder: "%author%",
		},
		{
			name:                          "Tags with nested fields",
			tagName:                       "author.1.books",
			expectedPositionalPlaceholder: "%author_1_books%",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPositionalPlaceholder, makePositionalPlaceholder(tc.tagName))
		})
	}
}
