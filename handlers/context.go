package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
)

var ErrNoSuchParam error = errors.New("no such param in the context")

const CurrentUser string = "currentUser"

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
	currentUser, err := GetFromContext[*models.User](CurrentUser, c)
	return currentUser, err
}

func GetIDFromContext(c *gin.Context) (uint, error) {

	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)

	if err != nil {
		return uint(0), fmt.Errorf("failed to parse uint: %w", err)
	}

	return uint(id), nil

}

func SetUser(c *gin.Context, value *models.User) {
	c.Set(CurrentUser, value)
}
