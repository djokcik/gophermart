package encrypt

import (
	"errors"
	"fmt"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

func CreateToken(secretKey string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer:    "gophermart",
		},
	})

	return token.SignedString([]byte(secretKey))
}

func ParseToken(accessToken string, secretKey string) (int, error) {
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
		return claims.Id, nil
	}

	return 0, model.ErrInvalidAccessToken
}

func GetJwtTokenByAuthHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("unauthorized")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", errors.New("unauthorized")
	}

	if headerParts[0] != "Bearer" {
		return "", errors.New("unauthorized")
	}

	return headerParts[1], nil
}
