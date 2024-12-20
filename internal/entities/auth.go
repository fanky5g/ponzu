package entities

import "time"

type CredentialType string

var CredentialTypePassword CredentialType = "password"

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
	return "credential_hashes" 
}

func (*CredentialHash) EntityName() string {
	return "CredentialHash"
}

type AuthToken struct {
	Expires time.Time
	Token   string
}

type RecoveryKey struct {
	Email string `json:"email"`
	Value string `json:"value"`
}

func (*RecoveryKey) GetRepositoryToken() string {
	return "recovery_keys" 
}

func (*RecoveryKey) EntityName() string {
	return "RecoveryKey"
}
