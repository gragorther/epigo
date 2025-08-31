package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const typeEmailVerification = "emailVerification"

type emailClaims struct {
	claims
	Email string `json:"email,omitzero"`
}

func CreateEmailVerification(jwtSecret []byte, userEmail string, audience string, issuer string) (token string, err error) {
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, emailClaims{Email: userEmail,
		claims: claims{RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			Audience: []string{audience},
			Issuer:   issuer,
			IssuedAt: jwt.NewNumericDate(time.Now()),
		}, Type: typeEmailVerification}})

	return generatedToken.SignedString(jwtSecret)
}

var ErrInvalidToken error = errors.New("invalid token")
var ErrInvalidEmail error = errors.New("invalid user email address")
var ErrEmptyEmailClaim error = errors.New("empty email claim")

func ParseEmailVerification(jwtSecret []byte, tokenString string, audience string, issuer string) (userEmail string, err error) {
	var claims emailClaims
	if err := parseToken(jwtSecret, tokenString, typeEmailVerification, issuer, []string{audience}, "", &claims); err != nil {
		return "", err
	}

	return claims.Email, nil
}
