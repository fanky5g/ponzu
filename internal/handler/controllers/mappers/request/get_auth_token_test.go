package request

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type GetAuthTokenTestSuite struct {
	suite.Suite
}

func (suite *GetAuthTokenTestSuite) TestGetAuthTokenFromAuthorizationHeader() {
	token := "jwt-token"

	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	mappedToken := GetAuthToken(req)
	assert.Equal(suite.T(), mappedToken, token)
}

func (suite *GetAuthTokenTestSuite) TestGetAuthTokenFromCookie() {
	token := "jwt-token"

	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:    "_token",
		Value:   token,
		Expires: time.Now().AddDate(0, 0, 1),
		Path:    "/",
	})

	mappedToken := GetAuthToken(req)
	assert.Equal(suite.T(), mappedToken, token)
}

func TestGetAuthToken(t *testing.T) {
	suite.Run(t, new(GetAuthTokenTestSuite))
}
