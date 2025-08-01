package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
)

var ErrNoSuchParam error = errors.New("no such param in the context")

func GetFromContext[T any](key string, c *gin.Context) (T, error) {
	// empty var
	var t T

	value, exists := c.Get(key)
	if !exists {
		return t, ErrNoSuchParam
	}
	typedValue := value.(T)
	return typedValue, nil
}

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	currentUser, err := GetFromContext[*models.User]("currentUser", c)
	return currentUser, err
}
