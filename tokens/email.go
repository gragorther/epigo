package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TypeEmailVerification = "emailVerification"

type EmailClaims struct {
	Claims
	Email string `json:"email,omitzero"`
}
type CreateEmailVerificationFunc func(userEmail string) (token string, err error)

func CreateEmailVerification(jwtSecret []byte, audience string, issuer string) CreateEmailVerificationFunc {
	return func(userEmail string) (token string, err error) {
		return createToken(jwtSecret, EmailClaims{Email: userEmail,
			Claims: Claims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
				Audience: []string{audience},
				Issuer:   issuer,
				IssuedAt: jwt.NewNumericDate(time.Now()),
			}, Type: TypeEmailVerification}})
	}

}

var ErrInvalidToken error = errors.New("invalid token")
var ErrInvalidEmail error = errors.New("invalid user email address")
var ErrEmptyEmailClaim error = errors.New("empty email claim")

type ParseEmailVerificationFunc func(tokenString string) (userEmail string, err error)

func ParseEmailVerification(jwtSecret []byte, audience string, issuer string) ParseEmailVerificationFunc {
	return func(tokenString string) (userEmail string, err error) {
		var claims EmailClaims
		if err := parseToken(jwtSecret, tokenString, TypeEmailVerification, []string{audience}, issuer, "", &claims); err != nil {
			return "", err
		}

		return claims.Email, nil
	}

}
