package response

import (
	"time"

	"github.com/fanky5g/ponzu/internal/auth"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/response"
)

func MapAuthTokenResponse(authToken *auth.AuthToken) *response.AuthTokenResponse {
	return &response.AuthTokenResponse{
		Expires: authToken.Expires.Format(time.RFC3339),
		Token:   authToken.Token,
	}
}
