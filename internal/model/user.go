package model

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

var ErrInvalidAccessToken = errors.New("invalid auth token")

type (
	Claims struct {
		jwt.StandardClaims
		Id int
	}

	UserRequestDto struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	UserResponseDto struct {
		Token string `json:"token"`
	}

	User struct {
		Id        int       `json:"id"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"createdAt"`
		Balance   Accrual   `json:"balance"`

		Password string
	}

	UserBalance struct {
		Current   Accrual `json:"current"`
		Withdrawn int     `json:"withdrawn"`
	}
)

func (u User) Validate() error {
	if len(u.Username) < 3 || len(u.Username) > 20 {
		return ErrUsernameLength
	}

	if u.Password == "" {
		return ErrPasswordEmpty
	}

	if len(u.Password) < 3 || len(u.Password) > 256 {
		return ErrPasswordLength
	}

	return nil
}

var (
	ErrUsernameLength = errors.New("validate username: invalid length")
	ErrPasswordEmpty  = errors.New("validate password: empty")
	ErrPasswordLength = errors.New("validate password: invalid length")
)
