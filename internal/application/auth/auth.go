package auth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fanky5g/ponzu/database"
	"github.com/fanky5g/ponzu/internal/entities"
	"github.com/nilslice/jwt"
)

type Service struct {
	users        database.Repository
	credentials  database.Repository
	recoveryKeys database.Repository
	config       database.Repository
}

func (s *Service) IsTokenValid(token string) (bool, error) {
	return jwt.Passes(token), nil
}

func (s *Service) getUserByEmail(email string) (*entities.User, error) {
	u, err := s.users.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	return u.(*entities.User), nil
}

func (s *Service) GetUserFromAuthToken(token string) (*entities.User, error) {
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

func (s *Service) NewToken(user *entities.User) (*entities.AuthToken, error) {
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

func (s *Service) GetRecoveryKey(email string) (*entities.RecoveryKey, error) {
	r, err := s.recoveryKeys.FindOneBy(map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	return r.(*entities.RecoveryKey), nil
}

func (s *Service) SetRecoveryKey(email string) (*entities.RecoveryKey, error) {
	recoveryKey, err := s.recoveryKeys.Insert(&entities.RecoveryKey{
		Email: email,
		Value: fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().Unix())).Int63()),
	})

	if err != nil {
		return nil, err
	}

	return recoveryKey.(*entities.RecoveryKey), nil
}

func New(
	config database.Repository,
	users database.Repository,
	credentials database.Repository,
	recoveryKeys database.Repository,
) (*Service, error) {
	c, err := config.Latest()
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

	return &Service{
		users:        users,
		credentials:  credentials,
		recoveryKeys: recoveryKeys,
		config:       config,
	}, nil
}
