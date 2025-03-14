package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fanky5g/ponzu/internal/database"
	"github.com/nilslice/jwt"
)

type DomainConfigInterface interface {
	GetDomain() (string, error)
}

type Service struct {
	userRepository        database.Repository
	credentialRepository  database.Repository
	recoveryKeyRepository database.Repository
	domainConfig          DomainConfigInterface
}

func (s *Service) IsTokenValid(token string) (bool, error) {
	return jwt.Passes(token), nil
}

func (s *Service) getUserByEmail(email string) (*User, error) {
	u, err := s.userRepository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*User), nil
}

func (s *Service) GetUserFromAuthToken(token string) (*User, error) {
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

func (s *Service) NewToken(user *User) (*Token, error) {
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

	return &Token{
		Expires: expires,
		Token:   token,
	}, nil
}

func (s *Service) GetRecoveryKey(email string) (*RecoveryKey, error) {
	r, err := s.recoveryKeyRepository.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	return r.(*RecoveryKey), nil
}

func (s *Service) SetRecoveryKey(email string) (*RecoveryKey, error) {
	recoveryKey, err := s.recoveryKeyRepository.Insert(&RecoveryKey{
		Email: email,
		Value: fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().Unix())).Int63()),
	})

	if err != nil {
		return nil, err
	}

	return recoveryKey.(*RecoveryKey), nil
}

func New(clientSecret string, db database.Database) (*Service, error) {
	if clientSecret != "" {
		jwt.Secret([]byte(clientSecret))
	}

	return &Service{
		userRepository:        db.GetRepositoryByToken(UserRepositoryToken),
		credentialRepository:  db.GetRepositoryByToken(CredentialHashRepositoryToken),
		recoveryKeyRepository: db.GetRepositoryByToken(RecoveryKeyRepositoryToken),
	}, nil
}
