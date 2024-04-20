package request

import (
	"bytes"
	"encoding/json"
	"github.com/fanky5g/ponzu/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"mime/multipart"
	"net/http"
	"testing"
)

type AuthRequestTestSuite struct {
	suite.Suite
}

func (suite *AuthRequestTestSuite) TestMapJSONRequest() {
	password := "secret"
	accountId := "john@doe.com"
	request := map[string]interface{}{
		"account_id": accountId,
		"credential": map[string]string{
			"type":  "password",
			"value": password,
		},
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(request); err != nil {
		suite.FailNow(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")

	expectedCredential := &entities.Credential{
		Type:  entities.CredentialTypePassword,
		Value: password,
	}

	mappedAccountId, mappedCredential, err := MapAuthRequest(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedCredential, mappedCredential)
		assert.Equal(suite.T(), accountId, mappedAccountId)
	}
}

func (suite *AuthRequestTestSuite) TestMapPostFormRequest() {
	password := "secret"
	accountId := "john@doe.com"

	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)

	w, err := multipartWriter.CreateFormField("account_id")
	if err != nil {
		suite.FailNow(err.Error())
	}

	if _, err = w.Write([]byte(accountId)); err != nil {
		suite.FailNow(err.Error())
	}

	w, err = multipartWriter.CreateFormField("credential_type")
	if err != nil {
		suite.FailNow(err.Error())
	}

	if _, err = w.Write([]byte("password")); err != nil {
		suite.FailNow(err.Error())
	}

	w, err = multipartWriter.CreateFormField("credential")
	if err != nil {
		suite.FailNow(err.Error())
	}

	if _, err = w.Write([]byte(password)); err != nil {
		suite.FailNow(err.Error())
	}

	if err = multipartWriter.Close(); err != nil {
		suite.FailNow(err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	expectedCredential := &entities.Credential{
		Type:  entities.CredentialTypePassword,
		Value: password,
	}

	mappedAccountId, mappedCredential, err := MapAuthRequest(req)
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), expectedCredential, mappedCredential)
		assert.Equal(suite.T(), accountId, mappedAccountId)
	}
}

func TestMapAuthRequest(t *testing.T) {
	suite.Run(t, new(AuthRequestTestSuite))
}
