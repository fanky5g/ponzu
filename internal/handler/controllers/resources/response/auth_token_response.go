package response

type AuthTokenResponse struct {
	Expires string `json:"expires"`
	Token   string `json:"token"`
}
