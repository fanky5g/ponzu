package request

import (
	"errors"
	"github.com/fanky5g/ponzu/internal/auth"
)

var (
	ErrInvalidAccountID      = errors.New("invalid account_id. account_id must be a non-empty string")
	ErrInvalidCredentialType = errors.New("invalid credential_type")
	ErrInvalidCredential     = errors.New("invalid credential")
	ErrUnsupportedAuthMethod = errors.New("unsupported auth method")
)

type Credential struct {
	Type  auth.CredentialType `json:"type"`
	Value interface{}         `json:"value"`
}

type AuthRequestDto struct {
	AccountID  string     `json:"account_id"`
	Credential Credential `json:"credential"`
}

func (request *AuthRequestDto) Validate() error {
	if request.AccountID == "" {
		return ErrInvalidAccountID
	}

	if request.Credential.Type == "" {
		return ErrInvalidCredentialType
	}

	if request.Credential.Value == "" {
		return ErrInvalidCredential
	}

	return nil
}

func (request *AuthRequestDto) ToCredential() (*auth.Credential, error) {
	switch request.Credential.Type {
	case auth.CredentialTypePassword:
		password, ok := request.Credential.Value.(string)
		if !ok {
			return nil, ErrInvalidCredential
		}

		return &auth.Credential{
			Type:  auth.CredentialTypePassword,
			Value: password,
		}, nil
	default:
		return nil, ErrUnsupportedAuthMethod
	}
}
