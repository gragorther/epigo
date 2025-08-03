package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gragorther/epigo/handlers"
	"github.com/gragorther/epigo/models"
	"github.com/stretchr/testify/assert"
)

func (m *mockDB) CreateLastMessage(lastMessage *models.LastMessage) error {
	m.LastMessages = append(m.LastMessages, *lastMessage)
	return m.Err
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
		messageInput := &handlers.MessageInput{
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
}
