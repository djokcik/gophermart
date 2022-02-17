package service

import "errors"

var (
	ErrWrongPassword = errors.New("authenticate: invalid username or password")
)
