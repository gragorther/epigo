package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TypeUserLifeStatus = "userLifeStatus"

func CreateUserLifeStatus(jwtSecret []byte, userID uint, expiresAfter time.Duration) (token string, err error) {
	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(expiresAfter).Unix(),
		"id":  userID,
		"typ": TypeUserLifeStatus,
	})
	return generatedToken.SignedString(jwtSecret)
}

/*
func ParseUserLifeStatus(jwtSecret []byte, token string) (userID uint, err error) {
	claims, err := parseToken(jwtSecret, token, TypeUserLifeStatus)

}
*/
