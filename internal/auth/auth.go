package auth

import (
	"time"
)

const (
	CredentialHashRepositoryToken = "credential_hashes"
	RecoveryKeyRepositoryToken    = "recovery_keys"
)

var CredentialTypePassword CredentialType = "password"

type CredentialType string

type PasswordHash struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

type Credential struct {
	Type  CredentialType `json:"type"`
	Value string         `json:"value"`
}

type CredentialHash struct {
	UserId string         `json:"user_id"`
	Type   CredentialType `json:"type"`
	Value  []byte         `json:"value"`
}

func (*CredentialHash) GetRepositoryToken() string {
	return CredentialHashRepositoryToken
}

func (*CredentialHash) EntityName() string {
	return "CredentialHash"
}

type Token struct {
	Expires time.Time
	Token   string
}

type RecoveryKey struct {
	Email string `json:"email"`
	Value string `json:"value"`
}

func (*RecoveryKey) GetRepositoryToken() string {
	return RecoveryKeyRepositoryToken
}

func (*RecoveryKey) EntityName() string {
	return "RecoveryKey"
}
