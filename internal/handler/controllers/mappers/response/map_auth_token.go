package response

import (
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/response"
	"time"
)

func MapAuthTokenResponse(authToken *entities.AuthToken) *response.AuthTokenResponse {
	return &response.AuthTokenResponse{
		Expires: authToken.Expires.Format(time.RFC3339),
		Token:   authToken.Token,
	}
}
