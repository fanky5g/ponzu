package request

import (
	"net/http"
	"strings"
)

func GetAuthToken(req *http.Request) string {
	if cookieToken := getAuthTokenFromCookie(req); cookieToken != "" {
		return cookieToken
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, "Bearer ")
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func getAuthTokenFromCookie(req *http.Request) string {
	// check if token exists in cookie
	cookie, err := req.Cookie("_token")
	if err != nil {
		return ""
	}

	return cookie.Value
}
