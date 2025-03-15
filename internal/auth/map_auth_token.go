package auth

import (
	"time"
)

func MapAuthTokenResponse(authToken *Token) *TokenResponse {
	return &TokenResponse{
		Expires: authToken.Expires.Format(time.RFC3339),
		Token:   authToken.Token,
	}
}
