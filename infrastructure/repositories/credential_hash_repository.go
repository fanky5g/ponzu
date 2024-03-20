package repositories

import (
	"github.com/fanky5g/ponzu/entities"
)

type CredentialHashRepositoryInterface interface {
	GetByUserId(userId string, credentialType entities.CredentialType) (*entities.CredentialHash, error)
	SetCredential(hash *entities.CredentialHash) error
}
