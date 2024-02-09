package request

import (
	"errors"
	"github.com/fanky5g/ponzu/internal/domain/entities"
)

var (
	ErrInvalidAccountID      = errors.New("invalid account_id. account_id must be a non-empty string")
	ErrInvalidCredentialType = errors.New("invalid credential_type")
	ErrInvalidCredential     = errors.New("invalid credential")
	ErrUnsupportedAuthMethod = errors.New("unsupported auth method")
)

type Credential struct {
	Type  entities.CredentialType `json:"type"`
	Value interface{}             `json:"value"`
}

type AuthRequest struct {
	AccountID  string     `json:"account_id"`
	Credential Credential `json:"credential"`
}

func (auth *AuthRequest) Validate() error {
	if auth.AccountID == "" {
		return ErrInvalidAccountID
	}

	if auth.Credential.Type == "" {
		return ErrInvalidCredentialType
	}

	if auth.Credential.Value == "" {
		return ErrInvalidCredential
	}

	return nil
}

func (auth *AuthRequest) ToCredential() (*entities.Credential, error) {
	switch auth.Credential.Type {
	case entities.CredentialTypePassword:
		password, ok := auth.Credential.Value.(string)
		if !ok {
			return nil, ErrInvalidCredential
		}

		return &entities.Credential{
			Type:  entities.CredentialTypePassword,
			Value: password,
		}, nil
	default:
		return nil, ErrUnsupportedAuthMethod
	}
}
