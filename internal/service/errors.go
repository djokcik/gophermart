package service

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrWrongPassword = errors.New("authenticate: invalid username or password")

	ErrNotAuthenticated                = errors.New("service: no authenticted user found in the context")
	ErrOrderAlreadyUploadedAnotherUser = errors.New("service: order already uploaded another user")
	ErrOrderAlreadyUploaded            = errors.New("service: order already uploaded")

	ErrInsufficientFunds = errors.New("service: insufficient funds")
)
