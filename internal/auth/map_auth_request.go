package auth

import (
	"encoding/json"
	"errors"
	"github.com/fanky5g/ponzu/internal/http/request"
	"net/http"
	"net/url"
)

var ErrInvalidRequest = errors.New("invalid request")

func MapAuthRequest(req *http.Request) (string, *Credential, error) {
	var authRequest *RequestDto
	contentType := request.GetContentType(req)

	if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
		if err := req.ParseMultipartForm(1024 * 1024 * 4); err != nil {
			return "", nil, err
		}

		authRequest = mapAuthRequestFromFormData(req.PostForm)
	} else if contentType == "application/json" {
		if err := json.NewDecoder(req.Body).Decode(&authRequest); err != nil {
			return "", nil, err
		}
	}

	if authRequest == nil {
		return "", nil, ErrInvalidRequest
	}

	if err := authRequest.Validate(); err != nil {
		return "", nil, err
	}

	credential, err := authRequest.ToCredential()
	if err != nil {
		return "", nil, err
	}

	return authRequest.AccountID, credential, nil
}

func mapAuthRequestFromFormData(values url.Values) *RequestDto {
	return &RequestDto{
		AccountID: values.Get("account_id"),
		Credential: CredentialRequest{
			Type:  CredentialType(values.Get("credential_type")),
			Value: values.Get("credential"),
		},
	}
}
