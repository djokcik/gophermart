package service

import (
	"fmt"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

//go:generate mockery --name=AuthService

type AuthService interface {
	CreateToken(secretKey string, id int) (string, error)
	ParseToken(accessToken string, secretKey string) (int, error)
	GetJwtTokenByAuthHeader(authHeader string) (string, error)
	HashAndSalt(pwd string, pepper string) (string, error)
	CompareHashAndPassword(password string, hash string) error
}

func NewAuthService() AuthService {
	return &authService{}
}

type authService struct {
}

func (a authService) CompareHashAndPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a authService) HashAndSalt(pwd string, pepper string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd+pepper), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcryptPassword: %w", err)
	}

	return string(hash), nil
}

func (a authService) GetJwtTokenByAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrUnauthorized
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", ErrUnauthorized
	}

	if headerParts[0] != "Bearer" {
		return "", ErrUnauthorized
	}

	return headerParts[1], nil
}

func (a authService) CreateToken(secretKey string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    "gophermart",
		},
	})

	return token.SignedString([]byte(secretKey))
}

func (a authService) ParseToken(accessToken string, secretKey string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims.ID, nil
	}

	return 0, model.ErrInvalidAccessToken
}
