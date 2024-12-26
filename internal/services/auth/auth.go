package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/database"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/nilslice/jwt"
)

type service struct {
	userRepository        database.Repository
	credentialRepository  database.Repository
	recoveryKeyRepository database.Repository
	configRepository      database.Repository
}

type Service interface {
	IsTokenValid(token string) (bool, error)
	GetUserFromAuthToken(token string) (*auth.User, error)
	NewToken(user *auth.User) (*auth.AuthToken, error)
	SetCredential(userId string, credential *auth.Credential) error
	VerifyCredential(userId string, credential *auth.Credential) error
	LoginByEmail(email string, credential *auth.Credential) (*auth.AuthToken, error)
	GetRecoveryKey(email string) (*auth.RecoveryKey, error)
	SetRecoveryKey(email string) (*auth.RecoveryKey, error)
	SendPasswordRecoveryInstructions(email string) error
}

func (s *service) IsTokenValid(token string) (bool, error) {
	return jwt.Passes(token), nil
}

func (s *service) getUserByEmail(email string) (*auth.User, error) {
	u, err := s.userRepository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*auth.User), nil
}

func (s *service) GetUserFromAuthToken(token string) (*auth.User, error) {
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

func (s *service) NewToken(user *auth.User) (*auth.AuthToken, error) {
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

	return &auth.AuthToken{
		Expires: expires,
		Token:   token,
	}, nil
}

func (s *service) GetRecoveryKey(email string) (*auth.RecoveryKey, error) {
	r, err := s.recoveryKeyRepository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	return r.(*auth.RecoveryKey), nil
}

func (s *service) SetRecoveryKey(email string) (*auth.RecoveryKey, error) {
	recoveryKey, err := s.recoveryKeyRepository.Insert(&auth.RecoveryKey{
		Email: email,
		Value: fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().Unix())).Int63()),
	})

	if err != nil {
		return nil, err
	}

	return recoveryKey.(*auth.RecoveryKey), nil
}

func New(db database.Database) (Service, error) {
	configRepository := db.GetRepositoryByToken(tokens.ConfigRepositoryToken)
	c, err := configRepository.Latest()
	if err != nil {
		return nil, err
	}

	// TODO: update jwt secret whenever config.ClientSecret is updated
	var cfg config.Config
	if c != nil {
		cfg = *(c.(*config.Config))
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
