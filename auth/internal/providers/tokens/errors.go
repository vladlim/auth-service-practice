package tokens

import "errors"

var (
	ErrAccessGenerate  = errors.New("access token generation error")
	ErrRefreshGenerate = errors.New("refresh token generation error")
)
