package tokens

import "github.com/golang-jwt/jwt/v5"

type claims struct {
	jwt.RegisteredClaims
	Type string `json:"typ,omitempty"`
}

func (c claims) CheckType(expected string) (match bool) {
	return c.Type == expected
}
