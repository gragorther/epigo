package handlers_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Err error
}

func (m *mockDB) CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error {
	return m.Err

}

// sets up gin and returns a gin Context and an http response recorder
func setupGin() (*gin.Context, *httptest.ResponseRecorder) {

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestAddGroup(t *testing.T) {
	assert := assert.New(t)
	t.Run("with valid json input", func(t *testing.T) {
		c, w := setupGin()

		mock := mockDB{
			Err: nil,
		}
		jsonInput := handlers.GroupInput{
			RecipientEmails: []string{"gregor@gregtech.eu", "test@uwu.com"},
		}
		jsonString, _ := sonic.Marshal(&jsonInput)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonString))

		handler := handlers.AddGroup(&mock)
		handler(c)

		assert.Equal(200, w.Code, "response code should be 200")
	})

}
