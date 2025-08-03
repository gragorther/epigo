package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

func setGinHttpBody(c *gin.Context, buf []byte) {
	c.Request = &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(buf)),
	}
}

func (m *mockDB) CreateLastMessage(lastMessage *models.LastMessage) error {
	m.LastMessages = append(m.LastMessages, *lastMessage)
	return m.Err
}
func (m *mockDB) FindLastMessagesByUserID(userID uint) ([]models.LastMessage, error) {
	var output []models.LastMessage

	for _, message := range m.LastMessages {
		if message.UserID == userID {
			output = append(output, message)
		}
	}

	return output, m.Err
}
func (m *mockDB) UpdateLastMessage(newMessage *models.LastMessage) error {
	for _, message := range m.LastMessages {
		if message.ID == newMessage.ID {
			message = *newMessage
		}
	}
	return m.Err
}

func (m *mockDB) CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error) {
	for _, message := range m.LastMessages {
		if message.ID == messageID {
			if message.UserID == userID {
				return true, m.Err
			}
		}
	}
	return false, m.Err
}

type invalidMessageInput struct {
	Title    string `json:"titleasdf"`
	Content  string `json:"contentasdf"`
	GroupIDs []uint `json:"groupIDsasdf"`
}

func TestAddLastMessage(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid input", func(t *testing.T) {
		c, w := setupGin()
		userName := "testname"
		currentUser := &models.User{ID: 1, Name: &userName}
		c.Set("currentUser", currentUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = true

		// TODO: use struct composition instead of this weird duplication
		messageInput := &handlers.AddMessageInput{
			Title:    "uwu",
			Content:  "I hereby declare AAAAAAAAAAAAAA",
			GroupIDs: []uint{1, 2, 3, 5}}
		jsonString, _ := sonic.Marshal(messageInput)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonString)),
		}

		handler := handlers.AddLastMessage(mock)
		handler(c)

		//the field in the mock db we're interested in
		definedField := mock.LastMessages[0]
		assert.Equal(http.StatusOK, w.Code, "http status code should indicate success")
		assert.Equal(messageInput.Title, definedField.Title)
		assert.Equal(messageInput.Content, definedField.Content)
		for i, group := range definedField.Groups {
			assert.Equal(messageInput.GroupIDs[i], group.ID)
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		c, w := setupGin()
		userName := "testname"
		currentUser := &models.User{ID: 1, Name: &userName}
		c.Set("currentUser", currentUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = true

		// TODO: use struct composition instead of this weird duplication
		messageInput := &invalidMessageInput{
			Title:    "uwu",
			Content:  "I hereby declare AAAAAAAAAAAAAA",
			GroupIDs: []uint{1, 2, 3, 5}}
		jsonString, _ := sonic.Marshal(messageInput)
		c.Request = &http.Request{
			Body: io.NopCloser(bytes.NewBuffer(jsonString)),
		}

		handler := handlers.AddLastMessage(mock)
		handler(c)

		//the field in the mock db we're interested in
		assert.Equal(http.StatusUnprocessableEntity, w.Code, "http status code should not indicate success")
		assert.Equal([]models.LastMessage(nil), mock.LastMessages, "there should be no last messages created because the input was invalid")

	})
	t.Run("user does not own group to which the last message is being added", func(t *testing.T) {
		c, w := setupGin()
		userName := "testname"
		currentUser := &models.User{ID: 1, Name: &userName}
		c.Set("currentUser", currentUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = false // not authorized
		messageInput, err := sonic.Marshal(handlers.AddMessageInput{
			Title:    "uwu",
			Content:  "markdown input",
			GroupIDs: []uint{1, 2, 3, 4, 5, 6, 7, 8, 9},
		})
		if err != nil {
			t.Fatalf("sonic failed to marshal json, %v", err)
		}
		setGinHttpBody(c, messageInput)

		handler := handlers.AddLastMessage(mock)
		handler(c)

		assert.Equal(http.StatusUnauthorized, w.Code, "user should be unauthorized to perform this action")
		assert.Equal([]models.LastMessage(nil), mock.LastMessages)
	})
}
func TestListLastMessages(t *testing.T) {
	assert := assert.New(t)
	t.Run("valid input", func(t *testing.T) {
		c, w := setupGin()
		userName := "testname"
		currentUser := &models.User{ID: 1, Name: &userName}
		mock := newMockDB(nil)
		mock.IsAuthorized = true

		handlers.SetUser(c, currentUser)

		handler := handlers.ListLastMessages(mock)
		handler(c)

		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestEditLastMessage(t *testing.T) {
	assert := assert.New(t)

	t.Run("valid input", func(t *testing.T) {
		c, w := setupGin()
		userName := "testname"
		userID := uint(1)
		currentUser := &models.User{ID: userID, Name: &userName}
		handlers.SetUser(c, currentUser)
		mock := newMockDB(nil)
		mock.IsAuthorized = true
		c.AddParam("id", "1")

		mock.LastMessages = []models.LastMessage{
			{ID: 1, Title: "stuff", UserID: userID},
		}

		jsonInput, err := sonic.Marshal(handlers.EditMessageInput{
			Title:    "test title",
			Content:  "test content",
			GroupIDs: []uint{1, 2, 4, 5, 6, 7},
		})
		if err != nil {
			t.Fatalf("failed to bind json: %v", err)
		}
		setGinHttpBody(c, jsonInput)

		handler := handlers.EditLastMessage(mock)
		handler(c)
		t.Log(c.Errors)

		assert.Equal(http.StatusCreated, w.Code)
	})
}
