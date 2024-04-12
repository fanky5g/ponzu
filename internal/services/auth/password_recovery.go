package auth

import (
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/entities"
	domainErrors "github.com/fanky5g/ponzu/errors"
	emailer "github.com/nilslice/email"
	log "github.com/sirupsen/logrus"
)

func (s *service) SendPasswordRecoveryInstructions(email string) error {
	_, err := s.userRepository.GetUserByEmail(email)
	if errors.Is(err, domainErrors.ErrNoUserExists) {
		return errors.New("no user exists")
	}

	// create temporary key to verify user
	key, err := s.SetRecoveryKey(email)
	if err != nil {
		return fmt.Errorf("failed to set account recovery key: %v", err)
	}

	cfgIface, err := s.configRepository.Latest()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	if cfgIface == nil {
		return errors.New("failed to get config")
	}

	cfg := cfgIface.(*entities.Config)
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

`, email, cfg.Domain, key, cfg.Domain)

	msg := emailer.Message{
		To:      email,
		From:    fmt.Sprintf("ponzu@%s", cfg.Domain),
		Subject: fmt.Sprintf("Account Recovery [%s]", cfg.Domain),
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
