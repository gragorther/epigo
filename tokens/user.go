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

func CreateUserAuth(jwtSecret []byte, userID uint, issuer string, audience []string) (token string, err error) {
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

var ErrInvalidTokenType error = errors.New("invalid token type")

func ParseUserAuth(jwtSecret []byte, tokenString string, issuer string, audience []string) (userID uint, err error) {
	var claims userAuthClaims
	if err := parseToken(jwtSecret, tokenString, TypeUserAuth, issuer, audience, "", &claims); err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
