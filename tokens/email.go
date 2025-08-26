package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const typeEmailVerification = "email_verification"

func CreateEmailVerification(jwtSecret string, userID uint) (token string, err error) {
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"typ": typeEmailVerification,
	})
	return generatedToken.SignedString([]byte(jwtSecret))
}

var ErrInvalidToken error = errors.New("invalid token")

func ParseEmailVerification(jwtSecret string, tokenString string) (userID uint, err error) {
	claims, err := parseToken(jwtSecret, tokenString, typeEmailVerification)
	if err != nil {
		return 0, err
	}
	return getID(claims)

}
