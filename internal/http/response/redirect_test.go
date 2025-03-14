package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRedirectLocation(t *testing.T) {
	tt := []struct {
		name                     string
		publicPath               string
		target                   string
		expectedRedirectLocation string
	}{
		{
			name:                     "EmptyPaths",
			publicPath:               "",
			target:                   "",
			expectedRedirectLocation: "/",
		},
		{
			name:                     "WithTargetAndPublicPath",
			publicPath:               "content-manager",
			target:                   "/login",
			expectedRedirectLocation: "/content-manager/login",
		},
		{
			name:                     "WithNoPublicPath",
			publicPath:               "",
			target:                   "login",
			expectedRedirectLocation: "/login",
		},
		{
			name:                     "WithQueryString",
			publicPath:               "/content-manager",
			target:                   "/edit?contentType=Page&id=1",
			expectedRedirectLocation: "/content-manager/edit?contentType=Page&id=1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			location, err := getRedirectLocation(tc.publicPath, tc.target)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.expectedRedirectLocation, location)
			}
		})
	}
}
