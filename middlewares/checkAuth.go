package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gragorther/epigo/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	U db.Users
}

func (h *AuthMiddleware) CheckAuth(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
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
		log.Printf("Token error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claimsUserID, claimsBool := claims["id"].(uint)

	//checks if type assertion failed
	if !claimsBool {
		log.Print("type assertion failed during authorization")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := h.U.GetUserByID(claimsUserID)
	if err != nil {
		log.Print("failed to get user by ID during auth")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUser", user)

	c.Next()

}
