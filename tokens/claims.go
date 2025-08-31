package tokens

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Type string `json:"typ,omitempty"`
}

func (c Claims) CheckType(expected string) (match bool) {
	return c.Type == expected
}
