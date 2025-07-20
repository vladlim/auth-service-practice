package tokens

import "errors"

var (
	ErrAccessGenerate  = errors.New("access token generation error")
	ErrRefreshGenerate = errors.New("refresh token generation error")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenParse      = errors.New("token parse error")
	ErrInvalidRole     = errors.New("invalid role")
	ErrInvalidKey      = errors.New("invalid key")
)
