package request

import (
	"encoding/json"
	"errors"
	"github.com/fanky5g/ponzu/internal/domain/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/request"
	"net/http"
	"net/url"
	"strings"
)

var ErrInvalidRequest = errors.New("invalid request")

func MapAuthRequest(req *http.Request) (string, *entities.Credential, error) {
	var authRequest *request.AuthRequest
	contentType := req.Header.Get("Content-Type")
	if strings.Contains(contentType, ";") {
		contentType = strings.TrimSpace(contentType[:strings.Index(contentType, ";")])
	}

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

func mapAuthRequestFromFormData(values url.Values) *request.AuthRequest {
	return &request.AuthRequest{
		AccountID: values.Get("account_id"),
		Credential: request.Credential{
			Type:  entities.CredentialType(values.Get("credential_type")),
			Value: values.Get("credential"),
		},
	}
}
