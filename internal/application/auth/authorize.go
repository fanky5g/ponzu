package auth

import (
	"errors"
	"github.com/fanky5g/ponzu/internal/entities"
)

func (s *Service) Authorize(currentUserToken string, credential *entities.Credential) error {
	user, err := s.GetUserFromAuthToken(currentUserToken)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("invalid user")
	}

	return s.VerifyCredential(user.ID, credential)
}
