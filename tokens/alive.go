package tokens

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TypeUserLifeStatus = "userLifeStatus"

type UserLifeStatusClaims struct {
	Claims
	UserID uint `json:"userID,omitzero"`
}

func NewUserLifeStatusClaims(userID uint, audience []string, issuer string, expiresAt time.Time) UserLifeStatusClaims {
	return UserLifeStatusClaims{
		Claims: NewClaims(TypeUserLifeStatus, audience, issuer, jwt.NewNumericDate(expiresAt), nil, strconv.FormatUint(uint64(userID), 10)),
		UserID: userID,
	}

}

type CreateUserLifeStatusFunc func(userID uint, expiresAt time.Time) (token string, err error)

// token for verifying that the user is still alive (or something similar, i.e. not kidnapped)
func CreateUserLifeStatus(jwtSecret []byte, audience []string, issuer string) CreateUserLifeStatusFunc {
	return func(userID uint, expiresAt time.Time) (token string, err error) {
		return createToken(jwtSecret, NewUserLifeStatusClaims(userID, audience, issuer, expiresAt))
	}
}

type ParseUserLifeStatusFunc func(tokenString string) (userID uint, err error)

func ParseUserLifeStatus(jwtSecret []byte, audience []string, issuer string) ParseUserLifeStatusFunc {
	return func(tokenString string) (userID uint, err error) {
		var claims UserLifeStatusClaims
		if err := parseToken(jwtSecret, tokenString, TypeUserLifeStatus, audience, issuer, "", &claims); err != nil {
			return 0, fmt.Errorf("failed to parse token: %w", err)
		}
		return claims.UserID, nil
	}

}
