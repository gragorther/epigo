package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/tokens"

	"github.com/gin-gonic/gin"
)

func CheckAuth(db interface {
	GetUserByID(ctx context.Context, ID uint) (models.User, error)
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
		valid, userID, err := tokens.ParseUserAuth(jwtSecret, tokenString)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse user auth token: %w", err))
			return
		}
		if !valid {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		if userID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByID(c, userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.Set("currentUser", user)

		c.Next()
	}
}
