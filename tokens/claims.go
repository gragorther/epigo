package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Type string `json:"typ,omitempty"`
}

func (c Claims) CheckType(expected string) (match bool) {
	return c.Type == expected
}

// expiresIn is the time after time.Now() at which the token expires, e.g. 2 * time.Hour
func NewClaims(claimsType string, audience []string, issuer string, expiresAt *jwt.NumericDate, notBefore *jwt.NumericDate, subject string) Claims {
	return Claims{
		Type: claimsType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Audience:  audience,
			ExpiresAt: expiresAt,
			NotBefore: notBefore,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

}
