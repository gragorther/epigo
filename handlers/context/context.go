package ginctx

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/middlewares"
)

var (
	ErrNoSuchParam error = errors.New("no such param in the context")
	ErrInvalidType       = errors.New("invalid type")
)

func Get[T any](key string, c *gin.Context) (T, error) {
	// empty var
	var t T

	value, exists := c.Get(key)
	if !exists {
		return t, ErrNoSuchParam
	}
	typedValue, ok := value.(T)
	if !ok {
		return t, ErrInvalidType
	}
	return typedValue, nil
}

func GetUserID(c *gin.Context) (uint, error) {
	currentUser, err := Get[uint](middlewares.CurrentUser, c)
	return currentUser, err
}

func GetID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return uint(0), fmt.Errorf("failed to parse uint: %w", err)
	}

	return uint(id), nil
}

func SetUserID(c *gin.Context, id uint) {
	c.Set(middlewares.CurrentUser, id)
}
