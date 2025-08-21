package handlers_test

import (
	"strconv"
	"testing"

	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

func TestGetKeyFromContext(t *testing.T) {
	assert := assert.New(t)
	t.Run("key exists", func(t *testing.T) {
		c, _ := SetupGin()
		testName := "test"
		expected := &models.User{
			ID:      1,
			Profile: &models.Profile{Name: &testName}}
		c.Set("currentUser", expected)

		respUser, err := handlers.GetFromContext[*models.User]("currentUser", c)
		assert.Equal(nil, err, "there should be no error")
		assert.Equal(expected, respUser, "user retrieved from the function should match original user")
	})

	t.Run("key doesn't exist", func(t *testing.T) {
		c, _ := SetupGin()

		respUser, err := handlers.GetFromContext[*models.User]("currentUser", c)
		assert.Equal(handlers.ErrNoSuchParam, err, "there should be no error")
		assert.Empty(respUser, "user should be empty")
	})
}

func TestGetIDFromContext(t *testing.T) {
	assert := assert.New(t)
	t.Run("ID exists", func(t *testing.T) {
		c, _ := SetupGin()
		id := "121"
		c.AddParam("id", id)

		respID, err := handlers.GetIDFromContext(c)

		assert.Nil(err)
		assert.Equal(id, strconv.FormatUint(uint64(respID), 10))
	})
}
