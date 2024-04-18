package auth

import (
	"fmt"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/nilslice/jwt"
	"math/rand"
	"time"
)

type service struct {
	userRepository        driver.Repository
	credentialRepository  driver.Repository
	recoveryKeyRepository driver.Repository
	configRepository      driver.Repository
}

type Service interface {
	IsTokenValid(token string) (bool, error)
	GetUserFromAuthToken(token string) (*entities.User, error)
	NewToken(user *entities.User) (*entities.AuthToken, error)
	SetCredential(userId string, credential *entities.Credential) error
	VerifyCredential(userId string, credential *entities.Credential) error
	LoginByEmail(email string, credential *entities.Credential) (*entities.AuthToken, error)
	GetRecoveryKey(email string) (*entities.RecoveryKey, error)
	SetRecoveryKey(email string) (*entities.RecoveryKey, error)
	SendPasswordRecoveryInstructions(email string) error
}

func (s *service) IsTokenValid(token string) (bool, error) {
	return jwt.Passes(token), nil
}

func (s *service) getUserByEmail(email string) (*entities.User, error) {
	u, err := s.userRepository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*entities.User), nil
}

func (s *service) GetUserFromAuthToken(token string) (*entities.User, error) {
	isValid, err := s.IsTokenValid(token)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, fmt.Errorf("error. Invalid token")
	}

	claims := jwt.GetClaims(token)
	email, ok := claims["user"]
	if !ok {
		return nil, fmt.Errorf("error. No user data found in request token")
	}

	return s.getUserByEmail(email.(string))
}

func (s *service) NewToken(user *entities.User) (*entities.AuthToken, error) {
	// create new token
	expires := time.Now().Add(time.Hour * 24 * 7)
	claims := map[string]interface{}{
		"exp":  expires,
		"user": user.Email,
	}

	token, err := jwt.New(claims)
	if err != nil {
		return nil, err
	}

	return &entities.AuthToken{
		Expires: expires,
		Token:   token,
	}, nil
}

func (s *service) GetRecoveryKey(email string) (*entities.RecoveryKey, error) {
	r, err := s.recoveryKeyRepository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	return r.(*entities.RecoveryKey), nil
}

func (s *service) SetRecoveryKey(email string) (*entities.RecoveryKey, error) {
	recoveryKey, err := s.recoveryKeyRepository.Insert(&entities.RecoveryKey{
		Email: email,
		Value: fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().Unix())).Int63()),
	})

	if err != nil {
		return nil, err
	}

	return recoveryKey.(*entities.RecoveryKey), nil
}

func New(db driver.Database) (Service, error) {
	configRepository := db.GetRepositoryByToken(tokens.ConfigRepositoryToken)
	c, err := configRepository.Latest()
	if err != nil {
		return nil, err
	}

	// TODO: update jwt secret whenever config.ClientSecret is updated
	var cfg entities.Config
	if c != nil {
		cfg = *(c.(*entities.Config))
	}

	if cfg.ClientSecret != "" {
		jwt.Secret([]byte(cfg.ClientSecret))
	}

	return &service{
		userRepository:        db.GetRepositoryByToken(tokens.UserRepositoryToken),
		credentialRepository:  db.GetRepositoryByToken(tokens.CredentialHashRepositoryToken),
		recoveryKeyRepository: db.GetRepositoryByToken(tokens.RecoveryKeyRepositoryToken),
	}, nil
}
