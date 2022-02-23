package handler

import "errors"

var (
	ErrNotAuthenticated = errors.New("handler: no authenticted user found in the context")
)
