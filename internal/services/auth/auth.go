package auth

import (
	"fmt"
	"github.com/fanky5g/ponzu/driver"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/infrastructure/repositories"
	"github.com/fanky5g/ponzu/tokens"
	"github.com/nilslice/jwt"
	"math/rand"
	"time"
)

type service struct {
	userRepository        repositories.UserRepositoryInterface
	credentialRepository  repositories.CredentialHashRepositoryInterface
	recoveryKeyRepository repositories.RecoveryKeyRepositoryInterface
}

type Service interface {
	IsTokenValid(token string) (bool, error)
	GetUserFromAuthToken(token string) (*entities.User, error)
	NewToken(user *entities.User) (*entities.AuthToken, error)
	SetCredential(userId string, credential *entities.Credential) error
	VerifyCredential(userId string, credential *entities.Credential) error
	LoginByEmail(email string, credential *entities.Credential) (*entities.AuthToken, error)
	GetRecoveryKey(email string) (string, error)
	SetRecoveryKey(email string) (string, error)
}

func (s *service) IsTokenValid(token string) (bool, error) {
	return jwt.Passes(token), nil
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

	return s.userRepository.GetUserByEmail(email.(string))
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

func (s *service) GetRecoveryKey(email string) (string, error) {
	return s.recoveryKeyRepository.GetRecoveryKey(email)
}

func (s *service) SetRecoveryKey(email string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	key := fmt.Sprintf("%d", r.Int63())

	return key, s.recoveryKeyRepository.SetRecoveryKey(email, key)
}

func New(db driver.Database) (Service, error) {
	configRepository := db.Get(tokens.ConfigRepositoryToken).(repositories.ConfigRepositoryInterface)
	userRepository := db.Get(tokens.UserRepositoryToken).(repositories.UserRepositoryInterface)
	credentialRepository := db.Get(tokens.CredentialHashRepositoryToken).(repositories.CredentialHashRepositoryInterface)
	recoveryKeyRepository := db.Get(tokens.RecoveryKeyRepositoryToken).(repositories.RecoveryKeyRepositoryInterface)

	clientSecret := configRepository.Cache().GetByKey("client_secret").(string)
	if clientSecret != "" {
		jwt.Secret([]byte(clientSecret))
	}

	return &service{
		userRepository:        userRepository,
		credentialRepository:  credentialRepository,
		recoveryKeyRepository: recoveryKeyRepository,
	}, nil
}
