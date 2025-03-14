package auth

type TokenResponse struct {
	Expires string `json:"expires"`
	Token   string `json:"token"`
}
