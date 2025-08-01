package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gragorther/epigo/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckAuth(db interface {
	GetUserByID(ID uint) (*models.User, error)
}) gin.HandlerFunc {
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
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token error: %w", err))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid auth token format: %v", claims))
			return
		}

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
