package auth

import (
	"errors"
	"github.com/fanky5g/ponzu/internal/auth"
)

func (s *service) Authorize(currentUserToken string, credential *auth.Credential) error {
	user, err := s.GetUserFromAuthToken(currentUserToken)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("invalid user")
	}

	return s.VerifyCredential(user.ID, credential)
}
