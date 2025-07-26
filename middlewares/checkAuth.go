package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserStore interface {
	GetUserByID(ID uint) (*models.User, error)
}

type AuthMiddleware struct {
	u UserStore
}

func NewAuthMiddleware(u UserStore) *AuthMiddleware {
	return &AuthMiddleware{u: u}
}

func (h *AuthMiddleware) CheckAuth(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrMissingAuthHeader)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrInvalidAuthTokenFormat)
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Token error: %v", err)
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrInvalidToken)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrInvalidAuthTokenFormat)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrExpiredToken)
		return
	}

	user, err := h.u.GetUserByID(uint(claims["id"].(float64)))
	if err != nil {
		log.Print("failed to get user by ID during auth")
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrFailedToGetUserID)
		return
	}

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUser", user)

	c.Next()

}
