package tokens

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type properTypeClaims interface {
	jwt.Claims
	CheckType(expected string) (match bool)
}

var signingMethod = jwt.SigningMethodHS256

var (
	ErrInvalidExpirationDate = errors.New("invalid expiration date")
	ErrTokenIsNil            = errors.New("token is nil")
	ErrTokenClaimsAreNil     = errors.New("token claims are nil")
	ErrInvalidClaimsType     = errors.New("invalid claims type")
)

// claims MUST be a pointer so that it can be populated
func parseToken(jwtSecret []byte, tokenString string, expectedType string, audience []string, issuer string, subject string, claims properTypeClaims) error {
	val := reflect.ValueOf(claims)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return fmt.Errorf("claims must be a non-nil pointer")
	}

	tokenConstraints := []jwt.ParserOption{
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithAudience(audience...),
	}
	if subject != "" {
		tokenConstraints = append(tokenConstraints, jwt.WithSubject(subject))
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if token.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	}, tokenConstraints...)
	if err != nil {
		return err
	}
	if token == nil {
		return ErrTokenIsNil
	}
	if !token.Valid {
		return ErrInvalidToken
	}
	if !claims.CheckType(expectedType) {
		return ErrInvalidClaimsType
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if exp == nil {
		return ErrInvalidExpirationDate
	}

	if time.Now().After(exp.Time) {
		return ErrExpiredToken
	}

	return nil
}

func createToken(jwtSecret []byte, claims jwt.Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
