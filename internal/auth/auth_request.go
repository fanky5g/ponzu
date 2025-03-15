package auth

import (
	"errors"
)

var (
	ErrInvalidAccountID      = errors.New("invalid account_id. account_id must be a non-empty string")
	ErrInvalidCredentialType = errors.New("invalid credential_type")
	ErrInvalidCredential     = errors.New("invalid credential")
	ErrUnsupportedAuthMethod = errors.New("unsupported auth method")
)

type CredentialRequest struct {
	Type  CredentialType `json:"type"`
	Value interface{}    `json:"value"`
}

type RequestDto struct {
	AccountID  string            `json:"account_id"`
	Credential CredentialRequest `json:"credential"`
}

func (request *RequestDto) Validate() error {
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

func (request *RequestDto) ToCredential() (*Credential, error) {
	switch request.Credential.Type {
	case CredentialTypePassword:
		password, ok := request.Credential.Value.(string)
		if !ok {
			return nil, ErrInvalidCredential
		}

		return &Credential{
			Type:  CredentialTypePassword,
			Value: password,
		}, nil
	default:
		return nil, ErrUnsupportedAuthMethod
	}
}
