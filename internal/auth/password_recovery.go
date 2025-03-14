package auth

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/exceptions"
	emailer "github.com/nilslice/email"
	log "github.com/sirupsen/logrus"
)

func (s *Service) SendPasswordRecoveryInstructions(email string) error {
	_, err := s.getUserByEmail(email)
	if errors.Is(err, exceptions.ErrNoUserExists) {
		return errors.New("no user exists")
	}

	// create temporary key to verify user
	key, err := s.SetRecoveryKey(email)
	if err != nil {
		return fmt.Errorf("failed to set account recovery key: %v", err)
	}

	domain, err := s.domainConfig.GetDomain()
	if err != nil {
		return fmt.Errorf("failed to get domain: %v", err)
	}

	body := fmt.Sprintf(`
There has been an account recovery request made for the user with email:
%s

To recover your account, please go to https://%s/recover/key and enter 
this email address along with the following secret key:

%s

If you did not make the request, ignore this message and your password 
will remain as-is.


Thank you,
Ponzu CMS at %s

`, email, domain, key, domain)

	msg := emailer.Message{
		To:      email,
		From:    fmt.Sprintf("ponzu@%s", domain),
		Subject: fmt.Sprintf("Account Recovery [%s]", domain),
		Body:    body,
	}

	// TODO: queue in application worker goroutines
	go func() {
		err = msg.Send()
		if err != nil {
			log.Println("Failed to send message to:", msg.To, "about", msg.Subject, "Error:", err)
		}
	}()

	return nil
}
