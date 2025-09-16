package tokens

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidID error = errors.New("invalid user ID")

var ErrExpiredToken = errors.New("expired token")

const TypeUserAuth = "userAuth"

type userAuthClaims struct {
	Claims
	UserID uint `json:"id,omitzero"`
}

type CreateUserAuthFunc func(userID uint) (token string, err error)

func CreateUserAuth(jwtSecret []byte, audience []string, issuer string) CreateUserAuthFunc {
	return func(userID uint) (token string, err error) {
		return createToken(jwtSecret, userAuthClaims{
			UserID: userID,
			Claims: Claims{
				Type: TypeUserAuth,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    issuer,
					Audience:  audience,
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					Subject:   strconv.FormatUint(uint64(userID), 10),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		})
	}
}

var ErrInvalidTokenType error = errors.New("invalid token type")

type ParseUserAuthFunc func(tokenString string) (userID uint, err error)

func ParseUserAuth(jwtSecret []byte, audience []string, issuer string) ParseUserAuthFunc {
	return func(tokenString string) (userID uint, err error) {
		var claims userAuthClaims
		if err := parseToken(jwtSecret, tokenString, TypeUserAuth, audience, issuer, "", &claims); err != nil {
			return 0, err
		}

		return claims.UserID, nil
	}
}
