package web

import (
	"net/http"
	"net/url"
	"strings"
)

const bearerHeader = "Bearer "
const accessToken = "access_token"

func GetToken(r *http.Request) string {
	var authHeader = r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, bearerHeader) {
		return strings.TrimPrefix(authHeader, bearerHeader)
	}
	return r.URL.Query().Get(accessToken)
}

func GetTokenSocket(r http.Header, requets url.URL) string {
	var authHeader = r.Get("Authorization")
	if strings.HasPrefix(authHeader, bearerHeader) {
		return strings.TrimPrefix(authHeader, bearerHeader)
	}
	return requets.Query().Get(accessToken)
}

func GetTokenPublic(r *http.Request) string {
	var authHeader = r.Header.Get("public")
	if strings.HasPrefix(authHeader, bearerHeader) {
		return strings.TrimPrefix(authHeader, bearerHeader)
	}
	return ""
}
