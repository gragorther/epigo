package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const typeEmailVerification = "email_verification"
const emailClaim = "email"

func CreateEmailVerification(jwtSecret string, userEmail string) (token string, err error) {
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		emailClaim: userEmail,
		"exp":      time.Now().Add(2 * time.Hour).Unix(),
		"typ":      typeEmailVerification,
	})
	return generatedToken.SignedString([]byte(jwtSecret))
}

var ErrInvalidToken error = errors.New("invalid token")
var ErrInvalidEmail error = errors.New("invalid user email address")
var ErrEmptyEmailClaim error = errors.New("empty email claim")

func getEmail(claims jwt.MapClaims) (email string, err error) {
	email, ok := claims[emailClaim].(string)
	if !ok {
		return "", ErrInvalidEmail
	}
	if email == "" {
		return "", ErrEmptyEmailClaim
	}
	return email, nil
}

func ParseEmailVerification(jwtSecret string, tokenString string) (userEmail string, err error) {
	claims, err := parseToken(jwtSecret, tokenString, typeEmailVerification)
	if err != nil {
		return "", err
	}
	return getEmail(claims)
}
