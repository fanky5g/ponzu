package auth

import (
	"errors"
)

func (s *Service) Authorize(currentUserToken string, credential *Credential) error {
	user, err := s.GetUserFromAuthToken(currentUserToken)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("invalid user")
	}

	return s.VerifyCredential(user.ID, credential)
}
