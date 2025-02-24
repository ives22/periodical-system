package token

import "errors"

const (
	TOKEN_COOKIE_NAME  = "access_token"
	TOKEN_GIN_KEY_NAME = "access_token"
)

var (
	CookieNotFound = errors.New("Cookie not found")
)
