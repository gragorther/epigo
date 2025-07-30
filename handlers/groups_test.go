package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Err        error
	Groups     []models.Group
	Recipients []models.Recipient
}

func (m *mockDB) CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.Recipient) error {
	group.Recipients = *recipientEmails
	m.Groups = append(m.Groups, *group)
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
		username := "test"
		fakeUser := &models.User{ID: 1, Name: &username}
		c.Set("currentUser", fakeUser)

		mock := mockDB{
			Err: nil,
		}
		jsonInput := handlers.GroupInput{
			Recipients: []models.Recipient{
				{Email: "gregor@gregtech.eu"},
				{Email: "test@gregtech.eu"},
			},
			Name:        "test name",
			Description: "test description",
		}
		jsonString, _ := sonic.Marshal(&jsonInput)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonString)),
		}

		handler := handlers.AddGroup(&mock)
		handler(c)
		assert.Equal(http.StatusOK, w.Code, "status code should indicate success")
		assert.Equal(jsonInput.Name, mock.Groups[0].Name)
		assert.Equal(jsonInput.Description, mock.Groups[0].Description)
		assert.Equal(fakeUser.ID, mock.Groups[0].UserID)
		assert.Equal(jsonInput.Recipients, mock.Groups[0].Recipients)

	})
	t.Run("with invalid json input", func(t *testing.T) {
		c, _ := setupGin()
		username := "test"
		fakeUser := &models.User{ID: 1, Name: &username}
		c.Set("currentUser", fakeUser)

	})
}
