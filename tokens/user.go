package tokens

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateUserAuth(jwtSecret string, userID uint) (token string, err error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err = generateToken.SignedString([]byte(jwtSecret))
	return
}

func ParseUserAuth(jwtSecret string, tokenString string) (valid bool, userID uint, err error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return false, 0, err
	}

	if !token.Valid {
		return false, 0, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, 0, nil
	}

	// checks if the claim has expired
	if float64(time.Now().Unix()) > claims["exp"].(float64) {

		return false, 0, nil
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return false, 0, nil
	}
	return true, uint(id), nil
}
