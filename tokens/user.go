package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidID error = errors.New("invalid user ID")

func getID(claims jwt.MapClaims) (uint, error) {
	id, ok := claims["id"].(float64)
	if !ok {
		return 0, ErrInvalidID
	}
	return uint(id), nil
}

func parseToken(jwtSecret string, tokenString string, expectedType string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	if float64(time.Now().Unix()) > claims["exp"].(float64) {

		return nil, ErrExpiredToken
	}
	typ, ok := claims["typ"].(string)
	if !ok {
		return nil, ErrInvalidTokenType
	}
	if typ != expectedType {
		return nil, ErrInvalidTokenType
	}
	return claims, nil

}

var ErrExpiredToken error = errors.New("expired token")

const typeUserSession = "userSession"

func CreateUserAuth(jwtSecret string, userID uint) (token string, err error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"typ": typeUserSession,
	})

	token, err = generateToken.SignedString([]byte(jwtSecret))
	return
}

var ErrInvalidTokenType error = errors.New("invalid token type")

func ParseUserAuth(jwtSecret string, tokenString string) (userID uint, err error) {

	claims, err := parseToken(jwtSecret, tokenString, typeUserSession)
	if err != nil {
		return 0, err
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}
	return uint(id), nil
}
