package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gragorther/epigo/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(db interface {
	GetUserByID(ID uint) (*models.User, error)
}, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token error: %w", err))
			return
		}
		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid auth token format: %v", claims))
			return
		}

		// checks if the claim has expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByID(uint(claims["id"].(float64)))
		if err != nil {
			log.Print("failed to get user by ID during auth")
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to check auth: %w", err))
			return
		}

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("currentUser", user)

		c.Next()
	}
}
