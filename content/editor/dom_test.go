package editor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddClassName(t *testing.T) {
	attrs := make(map[string]string)
	addClassName(attrs, "my-new-class")
	assert.Equal(t, map[string]string{
		"class": "my-new-class",
	}, attrs)

	attrs = nil
	addClassName(attrs, "my-new-class")
	assert.Equal(t, map[string]string(nil), attrs)

	attrs = map[string]string{
		"class": "h1 pad-50",
	}
	addClassName(attrs, "my-new-class")
	assert.Equal(t, map[string]string{
		"class": "h1 pad-50 my-new-class",
	}, attrs)


    attrs = map[string]string{
        "class": "my-new-class h1",
    }

    addClassName(attrs, "my-new-class")
    assert.Equal(t, map[string]string{
        "class": "my-new-class h1",
    }, attrs)

    attrs = map[string]string{}
    addClassName(attrs, "")
    assert.Equal(t, map[string]string{}, attrs)
}
